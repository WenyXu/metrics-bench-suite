package partition_query_loader

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWeightedPlanSelectorIsDeterministicAndHonorsHotShare(t *testing.T) {
	plans := []QueryPlan{{RegionID: 1}, {RegionID: 2}, {RegionID: 3}}
	selector := NewWeightedPlanSelector(plans, 1, 1, rand.New(rand.NewSource(7)))
	for range 10 {
		plan, ok := selector.Select()
		if !ok || plan.RegionID != 1 {
			t.Fatalf("hot selection = %#v, %t", plan, ok)
		}
	}
	selector = NewWeightedPlanSelector(plans, 1, 0, rand.New(rand.NewSource(7)))
	for range 10 {
		plan, ok := selector.Select()
		if !ok || plan.RegionID == 1 {
			t.Fatalf("cold selection = %#v, %t", plan, ok)
		}
	}
}

func TestWorkloadDelayAddsDeterministicJitter(t *testing.T) {
	config := DefaultConfig()
	config.Period, config.Jitter = time.Second, 500*time.Millisecond
	first := WorkloadDelay(config, rand.New(rand.NewSource(9)))
	second := WorkloadDelay(config, rand.New(rand.NewSource(9)))
	if first != second || first < config.Period || first > config.Period+config.Jitter {
		t.Fatalf("delay = %s", first)
	}
}

func TestRunWorkloadLimitsConcurrencyAndAttributesMetrics(t *testing.T) {
	config := DefaultConfig()
	config.DryRun, config.Duration, config.Period, config.Concurrency = false, 35*time.Millisecond, time.Millisecond, 2
	var mu sync.Mutex
	active, maximum := 0, 0
	result, err := runWorkload(context.Background(), config, []QueryPlan{{RegionID: 42}}, rand.New(rand.NewSource(1)), func(context.Context, QueryPlan, []any) error {
		mu.Lock()
		active++
		if active > maximum {
			maximum = active
		}
		mu.Unlock()
		time.Sleep(4 * time.Millisecond)
		mu.Lock()
		active--
		mu.Unlock()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if maximum > config.Concurrency {
		t.Fatalf("maximum concurrency = %d", maximum)
	}
	if metrics := result.ByRegion[42]; metrics.RequestCount == 0 || metrics.ErrorCount != 0 || metrics.LatencyP95 <= 0 {
		t.Fatalf("unexpected metrics: %#v", metrics)
	}
}

func TestRunWorkloadReportsPerSamplingIntervalMetrics(t *testing.T) {
	config := DefaultConfig()
	config.DryRun, config.Duration, config.Period, config.StatsInterval, config.Concurrency = false, 20*time.Millisecond, time.Millisecond, 5*time.Millisecond, 1
	result, err := runWorkload(context.Background(), config, []QueryPlan{{RegionID: 42}}, nil, func(context.Context, QueryPlan, []any) error { return nil })
	if err != nil || len(result.IntervalMetrics) == 0 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	for _, interval := range result.IntervalMetrics {
		if interval[42].RequestCount == 0 && len(interval) != 0 {
			t.Fatalf("invalid interval metrics: %#v", interval)
		}
	}
}

func TestProfilesUseIndependentPlanCadence(t *testing.T) {
	count := func(profile string) uint64 {
		config := DefaultConfig()
		config.DryRun, config.Profile, config.Duration, config.Period, config.StatsInterval, config.Concurrency = false, profile, 15*time.Millisecond, 5*time.Millisecond, 5*time.Millisecond, 2
		result, err := runWorkload(context.Background(), config, []QueryPlan{{RegionID: 1}}, nil, func(context.Context, QueryPlan, []any) error { return nil })
		if err != nil {
			t.Fatal(err)
		}
		return result.ByRegion[1].RequestCount
	}
	if sustained, periodic := count("sustained"), count("periodic"); sustained == 0 || periodic == 0 {
		t.Fatalf("sustained=%d periodic=%d", sustained, periodic)
	}
}

func TestRunWorkloadCancelsRequests(t *testing.T) {
	config := DefaultConfig()
	config.DryRun, config.Duration, config.Period = false, 10*time.Millisecond, 5*time.Millisecond
	cancelled := make(chan struct{})
	var once sync.Once
	_, err := runWorkload(context.Background(), config, []QueryPlan{{RegionID: 1}}, nil, func(ctx context.Context, _ QueryPlan, _ []any) error {
		<-ctx.Done()
		once.Do(func() { close(cancelled) })
		return ctx.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("request was not cancelled")
	}
}

func TestRunWorkloadDryRunDoesNotExecute(t *testing.T) {
	config := DefaultConfig()
	config.DryRun = true
	called := false
	result, err := runWorkload(context.Background(), config, []QueryPlan{{RegionID: 1}}, nil, func(context.Context, QueryPlan, []any) error { called = true; return nil })
	if err != nil {
		t.Fatal(err)
	}
	if called || len(result.ByRegion) != 0 {
		t.Fatalf("dry run executed workload: %#v", result)
	}
}

func TestRunWorkloadUsesFixedPerPlanWindowAndLatestTimestamp(t *testing.T) {
	config := DefaultConfig()
	config.Profile, config.PeriodMin, config.PeriodMax = "periodic", 5*time.Millisecond, 5*time.Millisecond
	config.TimeRangeMin, config.TimeRangeMax, config.Concurrency = 2*time.Second, 2*time.Second, 1
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	latest := time.Date(2026, time.July, 13, 10, 0, 0, 0, time.UTC)
	lookups := 0
	var gotArgs []any
	var logs []string
	_, err := runWorkloadWithLookupAndLogger(ctx, config, []QueryPlan{{LogicalTable: "cpu", PhysicalTable: "cpu_physical", RegionID: 9, Partition: "p0"}}, rand.New(rand.NewSource(7)), func(_ context.Context, _ QueryPlan, args []any) error {
		gotArgs = args
		cancel()
		return nil
	}, func(context.Context, QueryPlan) (time.Time, bool, error) {
		lookups++
		return latest, true, nil
	}, func(format string, args ...any) { logs = append(logs, fmt.Sprintf(format, args...)) })
	if err != nil || lookups == 0 || len(gotArgs) != 2 || len(logs) < 2 {
		t.Fatalf("err=%v lookups=%d args=%#v logs=%v", err, lookups, gotArgs, logs)
	}
	if gotArgs[0] != latest.Add(-2*time.Second) || gotArgs[1] != latest {
		t.Fatalf("window args=%#v", gotArgs)
	}
	if !strings.Contains(logs[0], "logical_table=cpu") || !strings.Contains(logs[0], "period=5ms") || !strings.Contains(logs[0], "time_range=2s") {
		t.Fatalf("scheduled log=%q", logs[0])
	}
}

func TestPlanAssignmentsArePersistentSeededAndPreserveHotSkew(t *testing.T) {
	config := DefaultConfig()
	config.HotPartitions, config.HotShare = 1, 0.8
	config.PeriodMin, config.PeriodMax = 10*time.Millisecond, 90*time.Millisecond
	config.TimeRangeMin, config.TimeRangeMax = time.Second, 9*time.Second
	plans := []QueryPlan{{LogicalTable: "a", Partition: "p0", LeaderDatanode: "node-a"}, {LogicalTable: "b", Partition: "p1", LeaderDatanode: "node-a"}, {LogicalTable: "c", Partition: "p2", LeaderDatanode: "node-b"}}
	left := assignPlans(plans, config, rand.New(rand.NewSource(17)))
	right := assignPlans(plans, config, rand.New(rand.NewSource(17)))
	if len(left) != 3 || left[0].period != right[0].period || left[0].timeRange != right[0].timeRange || left[0].interval != right[0].interval {
		t.Fatalf("assignments are not persistent/deterministic: %#v %#v", left, right)
	}
	hotIndex := -1
	for index, item := range left {
		if item.hot {
			hotIndex = index
			break
		}
	}
	if hotIndex < 0 {
		t.Fatalf("hot/cold cadence skew missing: %#v", left)
	}
	for index, item := range left {
		if index != hotIndex && item.hot || index != hotIndex && item.interval <= left[hotIndex].interval {
			t.Fatalf("hot/cold cadence skew missing: %#v", left)
		}
	}
	for index, item := range left {
		if item.hot {
			for coldIndex, cold := range left {
				if coldIndex != index && cold.timeRange >= item.timeRange {
					t.Fatalf("hot window is not larger: %#v", left)
				}
			}
		}
	}
}

func TestHotPlansPreferOneCoLocatedDatanodeDeterministically(t *testing.T) {
	config := DefaultConfig()
	config.HotPartitions = 2
	plans := []QueryPlan{{RegionID: 1, LeaderDatanode: "node-a"}, {RegionID: 2, LeaderDatanode: "node-a"}, {RegionID: 3, LeaderDatanode: "node-a"}, {RegionID: 4, LeaderDatanode: "node-b"}}
	left := assignPlans(plans, config, rand.New(rand.NewSource(3)))
	right := assignPlans(plans, config, rand.New(rand.NewSource(3)))
	for index := range left {
		if left[index].hot != right[index].hot {
			t.Fatalf("selection is not deterministic: %#v %#v", left, right)
		}
		if left[index].hot && left[index].plan.LeaderDatanode != "node-a" {
			t.Fatalf("hot plan not co-located: %#v", left)
		}
	}
}
