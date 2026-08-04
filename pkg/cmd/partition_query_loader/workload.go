package partition_query_loader

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RegionRequestMetrics groups workload measurements by the region targeted by a plan.
type RegionRequestMetrics map[uint64]RequestMetrics

// WorkloadResult contains the request measurements collected while the workload ran.
type WorkloadResult struct {
	ByRegion        RegionRequestMetrics
	IntervalMetrics []RegionRequestMetrics
}

// WeightedPlanSelector deterministically selects hot and non-hot plans using a supplied random source.
type WeightedPlanSelector struct {
	hot      []QueryPlan
	cold     []QueryPlan
	hotShare float64
	random   *rand.Rand
}

// NewWeightedPlanSelector builds a selector whose first hotPartitions plans are hot.
func NewWeightedPlanSelector(plans []QueryPlan, hotPartitions int, hotShare float64, random *rand.Rand) *WeightedPlanSelector {
	if hotPartitions > len(plans) {
		hotPartitions = len(plans)
	}
	if hotPartitions < 0 {
		hotPartitions = 0
	}
	if random == nil {
		random = rand.New(rand.NewSource(1))
	}
	return &WeightedPlanSelector{
		hot:      append([]QueryPlan(nil), plans[:hotPartitions]...),
		cold:     append([]QueryPlan(nil), plans[hotPartitions:]...),
		hotShare: hotShare,
		random:   random,
	}
}

// Select returns a plan according to the configured hot-share, or false when no plans exist.
func (selector *WeightedPlanSelector) Select() (QueryPlan, bool) {
	if len(selector.hot) == 0 && len(selector.cold) == 0 {
		return QueryPlan{}, false
	}
	plans := selector.cold
	if len(selector.hot) > 0 && (len(selector.cold) == 0 || selector.random.Float64() < selector.hotShare) {
		plans = selector.hot
	}
	return plans[selector.random.Intn(len(plans))], true
}

// WorkloadDelay returns the delay before the next periodic request.
func WorkloadDelay(config Config, random *rand.Rand) time.Duration {
	delay := config.Period
	if config.Jitter > 0 {
		if random == nil {
			random = rand.New(rand.NewSource(1))
		}
		delay += time.Duration(random.Int63n(int64(config.Jitter) + 1))
	}
	return delay
}

// RunWorkload executes bounded partition-targeted plans until the duration or context ends.
func RunWorkload(ctx context.Context, db *sql.DB, config Config, plans []QueryPlan) (WorkloadResult, error) {
	return runWorkloadWithLookup(ctx, config, plans, nil, func(requestCtx context.Context, plan QueryPlan, args []any) error {
		rows, err := db.QueryContext(requestCtx, plan.SQL, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
		}
		return rows.Err()
	}, func(requestCtx context.Context, plan QueryPlan) (time.Time, bool, error) {
		var timestamp sql.NullTime
		row := db.QueryRowContext(requestCtx, plan.MaxTimestampSQL, plan.partitionArgs...)
		if err := row.Scan(&timestamp); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return time.Time{}, false, nil
			}
			return time.Time{}, false, err
		}
		return timestamp.Time, timestamp.Valid, nil
	})
}

type queryExecutor func(context.Context, QueryPlan, []any) error

func runWorkload(ctx context.Context, config Config, plans []QueryPlan, random *rand.Rand, execute queryExecutor) (WorkloadResult, error) {
	return runWorkloadWithLookupAndLogger(ctx, config, plans, random, execute, func(context.Context, QueryPlan) (time.Time, bool, error) { return time.Now(), true, nil }, log.Printf)
}

type maxTimestampLookup func(context.Context, QueryPlan) (time.Time, bool, error)

type planAssignment struct {
	plan      QueryPlan
	period    time.Duration
	interval  time.Duration
	timeRange time.Duration
	hot       bool
}

func runWorkloadWithLookup(ctx context.Context, config Config, plans []QueryPlan, random *rand.Rand, execute queryExecutor, lookup maxTimestampLookup) (WorkloadResult, error) {
	return runWorkloadWithLookupAndLogger(ctx, config, plans, random, execute, lookup, log.Printf)
}

func runWorkloadWithLookupAndLogger(ctx context.Context, config Config, plans []QueryPlan, random *rand.Rand, execute queryExecutor, lookup maxTimestampLookup, logf func(string, ...any)) (WorkloadResult, error) {
	result := WorkloadResult{ByRegion: make(RegionRequestMetrics)}
	if config.Period != defaultPeriod && config.PeriodMin == defaultPeriodMin && config.PeriodMax == defaultPeriodMax {
		config.PeriodMin, config.PeriodMax = config.Period, config.Period
	}
	if config.TimeRange != defaultTimeRange && config.TimeRangeMin == defaultTimeRange && config.TimeRangeMax == defaultTimeRangeMax {
		config.TimeRangeMin, config.TimeRangeMax = config.TimeRange, config.TimeRange
	}
	if config.DryRun || len(plans) == 0 {
		return result, nil
	}
	if err := config.Validate(); err != nil {
		return result, err
	}
	if random == nil {
		seed := int64(1)
		if config.RandomSeed != nil {
			seed = *config.RandomSeed
		}
		random = rand.New(rand.NewSource(seed))
	}
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if config.Duration > 0 {
		runCtx, cancel = context.WithTimeout(ctx, config.Duration)
		defer cancel()
	}
	assignments := assignPlans(plans, config, random)
	if logf == nil {
		logf = log.Printf
	}
	for _, item := range assignments {
		logf("partition_query_loader scheduled_plan logical_table=%s physical_table=%s region_id=%d partition=%s period=%s interval=%s time_range=%s hot=%t", item.plan.LogicalTable, item.plan.PhysicalTable, item.plan.RegionID, item.plan.Partition, item.period, item.interval, item.timeRange, item.hot)
	}
	logHotDatanodes(assignments, logf)
	metrics := newRegionMetrics()
	intervalMetrics := newRegionMetrics()
	semaphore := make(chan struct{}, config.Concurrency)
	statsTicker := time.NewTicker(config.StatsInterval)
	defer statsTicker.Stop()
	var intervals []RegionRequestMetrics
	var runners sync.WaitGroup
	for _, item := range assignments {
		runners.Add(1)
		go func(item planAssignment) {
			defer runners.Done()
			timer := time.NewTimer(0)
			defer timer.Stop()
			for {
				select {
				case <-runCtx.Done():
					return
				case <-timer.C:
					select {
					case semaphore <- struct{}{}:
						started := time.Now()
						requestCtx, requestCancel := context.WithTimeout(runCtx, item.interval)
						latest, ok, err := lookup(requestCtx, item.plan)
						if err == nil && ok {
							err = execute(requestCtx, item.plan, item.plan.Arguments(latest.Add(-item.timeRange), latest))
						}
						requestCancel()
						if ok || err != nil {
							metrics.record(item.plan.RegionID, time.Since(started), err != nil)
							intervalMetrics.record(item.plan.RegionID, time.Since(started), err != nil)
						}
						<-semaphore
					default:
					}
					timer.Reset(item.interval)
				}
			}
		}(item)
	}
	done := make(chan struct{})
	go func() { runners.Wait(); close(done) }()
	for {
		select {
		case <-runCtx.Done():
			<-done
			intervals = append(intervals, intervalMetrics.snapshotAndReset())
			return WorkloadResult{ByRegion: metrics.snapshot(), IntervalMetrics: intervals}, nil
		case <-done:
			intervals = append(intervals, intervalMetrics.snapshotAndReset())
			return WorkloadResult{ByRegion: metrics.snapshot(), IntervalMetrics: intervals}, nil
		case <-statsTicker.C:
			intervals = append(intervals, intervalMetrics.snapshotAndReset())
		}
	}
}

func assignPlans(plans []QueryPlan, config Config, random *rand.Rand) []planAssignment {
	if random == nil {
		seed := int64(1)
		if config.RandomSeed != nil {
			seed = *config.RandomSeed
		}
		random = rand.New(rand.NewSource(seed))
	}
	hotCount := config.HotPartitions
	if hotCount > len(plans) {
		hotCount = len(plans)
	}
	if hotCount < 0 {
		hotCount = 0
	}
	hotIndices := chooseHotPlanIndices(plans, hotCount, random)
	hot := make(map[int]bool, len(hotIndices))
	for _, index := range hotIndices {
		hot[index] = true
	}
	assignments := make([]planAssignment, len(plans))
	for i, plan := range plans {
		isHot := hot[i]
		periodMin, periodMax := config.PeriodMin, config.PeriodMax
		if hotCount > 0 && hotCount < len(plans) {
			if isHot {
				periodMax = hotPeriodMax(config)
			} else {
				periodMin = coldPeriodMin(config)
			}
		}
		period := randomDuration(random, periodMin, periodMax)
		interval := period
		var timeRange time.Duration
		if isHot {
			timeRange = randomDuration(random, hotWindowMin(config), config.TimeRangeMax)
		} else {
			timeRange = randomDuration(random, config.TimeRangeMin, coldWindowMax(config))
		}
		assignments[i] = planAssignment{plan: plan, period: period, interval: interval, timeRange: timeRange, hot: isHot}
	}
	return assignments
}

func chooseHotPlanIndices(plans []QueryPlan, hotCount int, random *rand.Rand) []int {
	if hotCount <= 0 || len(plans) == 0 {
		return nil
	}
	groups := make(map[string][]int)
	for index, plan := range plans {
		datanode := plan.LeaderDatanode
		if datanode == "" {
			datanode = "unknown"
		}
		groups[datanode] = append(groups[datanode], index)
	}
	bestSize := 0
	best := make([]string, 0)
	for datanode, indices := range groups {
		if len(indices) > bestSize {
			bestSize, best = len(indices), []string{datanode}
		} else if len(indices) == bestSize {
			best = append(best, datanode)
		}
	}
	sort.Strings(best)
	selectedDatanode := best[0]
	if len(best) > 1 {
		selectedDatanode = best[random.Intn(len(best))]
	}
	candidates := append([]int(nil), groups[selectedDatanode]...)
	shuffleInts(candidates, random)
	selected := append([]int(nil), candidates[:minInt(hotCount, len(candidates))]...)
	if len(selected) < hotCount {
		remaining := make([]int, 0, len(plans)-len(selected))
		chosen := make(map[int]bool, len(selected))
		for _, index := range selected {
			chosen[index] = true
		}
		for index := range plans {
			if !chosen[index] {
				remaining = append(remaining, index)
			}
		}
		shuffleInts(remaining, random)
		selected = append(selected, remaining[:hotCount-len(selected)]...)
	}
	return selected
}

func shuffleInts(values []int, random *rand.Rand) {
	for index := len(values) - 1; index > 0; index-- {
		swap := random.Intn(index + 1)
		values[index], values[swap] = values[swap], values[index]
	}
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func hotWindowMin(config Config) time.Duration {
	split := config.TimeRangeMin + (config.TimeRangeMax-config.TimeRangeMin)/2
	if split < config.TimeRangeMax {
		return split + time.Nanosecond
	}
	return split
}

func hotPeriodMax(config Config) time.Duration {
	share := config.HotShare
	if share < 0.5 {
		share = 0.5
	}
	return config.PeriodMin + time.Duration(float64(config.PeriodMax-config.PeriodMin)*(1-share))
}

func coldPeriodMin(config Config) time.Duration {
	minimum := hotPeriodMax(config)
	if minimum < config.PeriodMax {
		minimum++
	}
	return minimum
}

func coldWindowMax(config Config) time.Duration {
	return config.TimeRangeMin + (config.TimeRangeMax-config.TimeRangeMin)/2
}

func logHotDatanodes(assignments []planAssignment, logf func(string, ...any)) {
	regions := make(map[string]map[uint64]bool)
	for _, assignment := range assignments {
		if !assignment.hot {
			continue
		}
		datanode := assignment.plan.LeaderDatanode
		if datanode == "" {
			datanode = "unknown"
		}
		if regions[datanode] == nil {
			regions[datanode] = make(map[uint64]bool)
		}
		regions[datanode][assignment.plan.RegionID] = true
	}
	datanodes := make([]string, 0, len(regions))
	for datanode := range regions {
		datanodes = append(datanodes, datanode)
	}
	sort.Strings(datanodes)
	for _, datanode := range datanodes {
		ids := make([]uint64, 0, len(regions[datanode]))
		for id := range regions[datanode] {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		formatted := make([]string, len(ids))
		for i, id := range ids {
			formatted[i] = strconv.FormatUint(id, 10)
		}
		logf("partition_query_loader hot_datanode datanode_id=%s regions=%s", datanode, strings.Join(formatted, ","))
	}
}

func randomDuration(random *rand.Rand, min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	return min + time.Duration(random.Int63n(int64(max-min)+1))
}

type regionMetrics struct {
	mu        sync.Mutex
	counts    map[uint64]uint64
	errors    map[uint64]uint64
	latencies map[uint64][]time.Duration
}

func newRegionMetrics() *regionMetrics {
	return &regionMetrics{counts: make(map[uint64]uint64), errors: make(map[uint64]uint64), latencies: make(map[uint64][]time.Duration)}
}

func (metrics *regionMetrics) record(regionID uint64, latency time.Duration, failed bool) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.counts[regionID]++
	if failed {
		metrics.errors[regionID]++
	}
	metrics.latencies[regionID] = append(metrics.latencies[regionID], latency)
}

func (metrics *regionMetrics) snapshot() RegionRequestMetrics {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	result := make(RegionRequestMetrics, len(metrics.counts))
	for regionID, count := range metrics.counts {
		latencies := append([]time.Duration(nil), metrics.latencies[regionID]...)
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		p95 := time.Duration(0)
		if len(latencies) > 0 {
			p95 = latencies[(len(latencies)*95+99)/100-1]
		}
		result[regionID] = RequestMetrics{RequestCount: count, ErrorCount: metrics.errors[regionID], LatencyP95: p95}
	}
	return result
}

func (metrics *regionMetrics) snapshotAndReset() RegionRequestMetrics {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	result := metrics.snapshotLocked()
	metrics.counts = make(map[uint64]uint64)
	metrics.errors = make(map[uint64]uint64)
	metrics.latencies = make(map[uint64][]time.Duration)
	return result
}

func (metrics *regionMetrics) snapshotLocked() RegionRequestMetrics {
	result := make(RegionRequestMetrics, len(metrics.counts))
	for regionID, count := range metrics.counts {
		latencies := append([]time.Duration(nil), metrics.latencies[regionID]...)
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		p95 := time.Duration(0)
		if len(latencies) > 0 {
			p95 = latencies[(len(latencies)*95+99)/100-1]
		}
		result[regionID] = RequestMetrics{RequestCount: count, ErrorCount: metrics.errors[regionID], LatencyP95: p95}
	}
	return result
}
