package partition_query_loader

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type workflowDependencies struct {
	openDB              func(Config) (*sql.DB, error)
	discover            func(context.Context, *sql.DB, []string) (DiscoveryResult, error)
	discoverWithOptions func(context.Context, *sql.DB, []string, DiscoveryOptions) (DiscoveryResult, error)
	runWorkload         func(context.Context, *sql.DB, Config, []QueryPlan) (WorkloadResult, error)
	sample              func(context.Context, *sql.DB, []DiscoveredTable, time.Time) ([]RegionSnapshot, error)
	write               func(io.Writer, string, []Observation) error
	now                 func() time.Time
	logf                func(string, ...any)
}

func defaultWorkflowDependencies() workflowDependencies {
	return workflowDependencies{
		openDB:              func(config Config) (*sql.DB, error) { return config.OpenDB() },
		discover:            Discover,
		discoverWithOptions: DiscoverWithOptions,
		runWorkload:         RunWorkload,
		sample:              SampleRegionStatistics,
		write:               WriteObservations,
		now:                 time.Now,
		logf:                log.Printf,
	}
}

// RunFunc runs the loader after command-line configuration is validated.
type RunFunc func(Config) error

// NewCommand creates the partition query loader command.
func NewCommand() *cobra.Command {
	return newCommand(func(ctx context.Context, config Config) error {
		_, err := RunConfigured(ctx, config, os.Stdout)
		return err
	})
}

// RunConfigured discovers safe plans and runs or reports the configured read-only workflow.
func RunConfigured(ctx context.Context, config Config, writer io.Writer) (WorkloadResult, error) {
	return runConfigured(ctx, config, writer, defaultWorkflowDependencies())
}

func runConfigured(ctx context.Context, config Config, writer io.Writer, dependencies workflowDependencies) (WorkloadResult, error) {
	if err := config.Validate(); err != nil {
		return WorkloadResult{}, err
	}
	db, err := dependencies.openDB(config)
	if err != nil {
		return WorkloadResult{}, err
	}
	if db != nil {
		defer db.Close()
	}

	discoveryCtx, cancelDiscovery := context.WithTimeout(ctx, config.DiscoveryTimeout)
	var discovery DiscoveryResult
	if dependencies.discoverWithOptions != nil {
		discovery, err = dependencies.discoverWithOptions(discoveryCtx, db, config.Databases, DiscoveryOptions{TablesPerDatabase: config.TablesPerDatabase, RandomSeed: config.RandomSeed, Progress: func(database string, selected, total, completed int) {
			if config.DryRun {
				fmt.Fprintf(writer, "discovery database=%s selected=%d completed=%d/%d\n", database, selected, completed, total)
			}
		}})
	} else {
		discovery, err = dependencies.discover(discoveryCtx, db, config.Databases)
	}
	cancelDiscovery()
	if err != nil {
		return WorkloadResult{}, err
	}
	planned := BuildQueryPlans(discovery)
	if dependencies.logf == nil {
		dependencies.logf = log.Printf
	}
	duration := config.Duration.String()
	if config.Duration == 0 {
		duration = "unlimited"
	}
	dependencies.logf("partition_query_loader initialized eligible_logical_tables=%d query_plans=%d profile=%s concurrency=%d duration=%s dry_run=%t", len(discovery.Tables), len(planned.Plans), config.Profile, config.Concurrency, duration, config.DryRun)
	if config.DryRun {
		for _, plan := range planned.Plans {
			if _, err := fmt.Fprintf(writer, "region_id=%d partition=%s sql=%s\n", plan.RegionID, plan.Partition, plan.SQL); err != nil {
				return WorkloadResult{}, err
			}
		}
		for _, skipped := range planned.Skipped {
			if _, err := fmt.Fprintf(writer, "skipped database=%s table=%s reason=%s\n", skipped.Database, skipped.Table, skipped.Reason); err != nil {
				return WorkloadResult{}, err
			}
		}
		return WorkloadResult{ByRegion: make(RegionRequestMetrics)}, nil
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	previous, err := dependencies.sample(runCtx, db, discovery.Tables, dependencies.now())
	if err != nil {
		return WorkloadResult{}, err
	}
	snapshots := [][]RegionSnapshot{previous}
	logDatanodeCPUAggregates(previous, nil, dependencies.logf)
	workloadDone := make(chan struct {
		result WorkloadResult
		err    error
	}, 1)
	go func() {
		result, err := dependencies.runWorkload(runCtx, db, config, planned.Plans)
		workloadDone <- struct {
			result WorkloadResult
			err    error
		}{result, err}
	}()

	ticker := time.NewTicker(config.StatsInterval)
	defer ticker.Stop()
	for {
		select {
		case completed := <-workloadDone:
			if completed.err != nil {
				return completed.result, completed.err
			}
			current, err := dependencies.sample(runCtx, db, discovery.Tables, dependencies.now())
			if err != nil {
				return completed.result, err
			}
			snapshots = append(snapshots, current)
			logDatanodeCPUAggregates(current, snapshots[len(snapshots)-2], dependencies.logf)
			observations := observationsFromSnapshots(snapshots, completed.result.IntervalMetrics)
			if err := dependencies.write(writer, config.OutputFormat, observations); err != nil {
				return completed.result, err
			}
			return completed.result, nil
		case <-ticker.C:
			current, err := dependencies.sample(runCtx, db, discovery.Tables, dependencies.now())
			if err != nil {
				cancel()
				<-workloadDone
				return WorkloadResult{}, err
			}
			snapshots = append(snapshots, current)
			logDatanodeCPUAggregates(current, snapshots[len(snapshots)-2], dependencies.logf)
		case <-ctx.Done():
			cancel()
			<-workloadDone
			return WorkloadResult{}, ctx.Err()
		}
	}
}

func logDatanodeCPUAggregates(current, previous []RegionSnapshot, logf func(string, ...any)) {
	for _, stats := range AggregateDatanodeCPU(current, previous) {
		logf("partition_query_loader datanode_cpu datanode_id=%s regions=%d total_cores=%g", stats.DatanodeID, stats.RegionCount, stats.TotalCores)
	}
}

func observationsFromSnapshots(snapshots [][]RegionSnapshot, metrics []RegionRequestMetrics) []Observation {
	var observations []Observation
	var previous []RegionSnapshot
	for index, current := range snapshots {
		for _, snapshot := range uniqueRegionSnapshots(current) {
			observations = append(observations, CalculateObservations([]RegionSnapshot{snapshot}, snapshotsForRegion(previous, snapshot.RegionID), intervalMetrics(index, metrics, snapshot.RegionID))...)
		}
		for _, snapshot := range uniqueRegionSnapshots(previous) {
			if len(snapshotsForRegion(current, snapshot.RegionID)) == 0 {
				observations = append(observations, CalculateObservations(nil, []RegionSnapshot{snapshot}, intervalMetrics(index, metrics, snapshot.RegionID))...)
			}
		}
		previous = current
	}
	return observations
}

func snapshotsForRegion(snapshots []RegionSnapshot, regionID uint64) []RegionSnapshot {
	for _, snapshot := range snapshots {
		if snapshot.RegionID == regionID {
			return []RegionSnapshot{snapshot}
		}
	}
	return nil
}

func uniqueRegionSnapshots(snapshots []RegionSnapshot) []RegionSnapshot {
	seen := make(map[uint64]struct{}, len(snapshots))
	unique := make([]RegionSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		if _, ok := seen[snapshot.RegionID]; !ok {
			seen[snapshot.RegionID] = struct{}{}
			unique = append(unique, snapshot)
		}
	}
	return unique
}

func intervalMetrics(index int, metrics []RegionRequestMetrics, regionID uint64) RequestMetrics {
	if index == 0 || index-1 >= len(metrics) {
		return RequestMetrics{}
	}
	return metrics[index-1][regionID]
}

// NewCommandWithRunner creates the command with an explicit execution dependency.
func NewCommandWithRunner(run RunFunc) *cobra.Command {
	if run == nil {
		return newCommand(nil)
	}
	return newCommand(func(_ context.Context, config Config) error { return run(config) })
}

func newCommand(run func(context.Context, Config) error) *cobra.Command {
	config := DefaultConfig()
	command := &cobra.Command{
		Use:          "partition_query_loader",
		Short:        "Generate safe, partition-targeted GreptimeDB read load",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			loaderConfig, err := ConfigFromCommand(cmd)
			if err != nil {
				return err
			}
			if run == nil {
				return fmt.Errorf("partition query loader runner is not configured")
			}
			return run(cmd.Context(), loaderConfig)
		},
	}

	flags := command.Flags()
	flags.StringVar(&config.MySQLHost, "mysql-host", config.MySQLHost, "MySQL host")
	flags.IntVar(&config.MySQLPort, "mysql-port", config.MySQLPort, "MySQL port")
	flags.StringVar(&config.MySQLUser, "mysql-user", "", "MySQL username")
	flags.StringVar(&config.MySQLPassword, "mysql-password", "", "MySQL password")
	flags.String("databases", defaultDatabases, "Comma-separated databases to discover")
	flags.IntVar(&config.TablesPerDatabase, "tables-per-database", config.TablesPerDatabase, "Maximum eligible logical tables selected per database; zero selects all")
	flags.Int64("random-seed", 0, "Optional seed for reproducible random table selection")
	flags.IntVar(&config.Concurrency, "concurrency", config.Concurrency, "Maximum concurrent queries")
	flags.DurationVar(&config.Duration, "duration", config.Duration, "Total load generation duration; 0 runs until cancelled")
	flags.DurationVar(&config.Period, "period", config.Period, "Deprecated alias for period-min and period-max")
	flags.DurationVar(&config.PeriodMin, "period-min", config.PeriodMin, "Inclusive minimum per-plan execution period")
	flags.DurationVar(&config.PeriodMax, "period-max", config.PeriodMax, "Inclusive maximum per-plan execution period")
	flags.DurationVar(&config.Jitter, "jitter", config.Jitter, "Maximum random delay added to each period")
	flags.DurationVar(&config.TimeRange, "time-range", config.TimeRange, "Deprecated alias for time-range-min and time-range-max")
	flags.DurationVar(&config.TimeRangeMin, "time-range-min", config.TimeRangeMin, "Inclusive minimum per-plan time range")
	flags.DurationVar(&config.TimeRangeMax, "time-range-max", config.TimeRangeMax, "Inclusive maximum per-plan time range")
	flags.IntVar(&config.HotPartitions, "hot-partitions", config.HotPartitions, "Number of hot partitions")
	flags.Float64Var(&config.HotShare, "hot-share", config.HotShare, "Share of requests sent to hot partitions")
	flags.DurationVar(&config.StatsInterval, "stats-interval", config.StatsInterval, "Region statistics sampling interval")
	flags.StringVar(&config.Profile, "profile", config.Profile, "Load profile: sustained or periodic")
	flags.String("workload-config", "", "Optional YAML workload settings; explicit flags take precedence")
	flags.StringVar(&config.OutputFormat, "output-format", config.OutputFormat, "Observation format: ndjson or csv")
	flags.BoolVar(&config.DryRun, "dry-run", config.DryRun, "Discover and print planned queries without executing them")
	flags.DurationVar(&config.DiscoveryTimeout, "discovery-timeout", config.DiscoveryTimeout, "Timeout for metadata discovery")
	flags.DurationVar(&config.DialTimeout, "mysql-dial-timeout", config.DialTimeout, "MySQL connection timeout")
	flags.DurationVar(&config.ReadTimeout, "mysql-read-timeout", config.ReadTimeout, "MySQL read timeout")
	flags.DurationVar(&config.WriteTimeout, "mysql-write-timeout", config.WriteTimeout, "MySQL write timeout")

	return command
}

// ConfigFromCommand parses and validates the command's flags.
func ConfigFromCommand(command *cobra.Command) (Config, error) {
	config := DefaultConfig()
	var err error
	if config.MySQLHost, err = command.Flags().GetString("mysql-host"); err != nil {
		return Config{}, err
	}
	if config.MySQLPort, err = command.Flags().GetInt("mysql-port"); err != nil {
		return Config{}, err
	}
	if config.MySQLUser, err = command.Flags().GetString("mysql-user"); err != nil {
		return Config{}, err
	}
	if config.MySQLPassword, err = command.Flags().GetString("mysql-password"); err != nil {
		return Config{}, err
	}
	databases, err := command.Flags().GetString("databases")
	if err != nil {
		return Config{}, err
	}
	config.Databases = parseDatabases(databases)
	if command.Flags().Changed("random-seed") {
		seed, err := command.Flags().GetInt64("random-seed")
		if err != nil {
			return Config{}, err
		}
		config.RandomSeed = &seed
	}
	if config.TablesPerDatabase, err = command.Flags().GetInt("tables-per-database"); err != nil {
		return Config{}, err
	}
	if config.Concurrency, err = command.Flags().GetInt("concurrency"); err != nil {
		return Config{}, err
	}
	if config.Duration, err = command.Flags().GetDuration("duration"); err != nil {
		return Config{}, err
	}
	if config.Period, err = command.Flags().GetDuration("period"); err != nil {
		return Config{}, err
	}
	if config.PeriodMin, err = command.Flags().GetDuration("period-min"); err != nil {
		return Config{}, err
	}
	if config.PeriodMax, err = command.Flags().GetDuration("period-max"); err != nil {
		return Config{}, err
	}
	if command.Flags().Changed("period") && !command.Flags().Changed("period-min") && !command.Flags().Changed("period-max") {
		config.PeriodMin, config.PeriodMax = config.Period, config.Period
	}
	if config.Jitter, err = command.Flags().GetDuration("jitter"); err != nil {
		return Config{}, err
	}
	if config.TimeRange, err = command.Flags().GetDuration("time-range"); err != nil {
		return Config{}, err
	}
	if config.TimeRangeMin, err = command.Flags().GetDuration("time-range-min"); err != nil {
		return Config{}, err
	}
	if config.TimeRangeMax, err = command.Flags().GetDuration("time-range-max"); err != nil {
		return Config{}, err
	}
	if command.Flags().Changed("time-range") && !command.Flags().Changed("time-range-min") && !command.Flags().Changed("time-range-max") {
		config.TimeRangeMin, config.TimeRangeMax = config.TimeRange, config.TimeRange
	}
	if config.HotPartitions, err = command.Flags().GetInt("hot-partitions"); err != nil {
		return Config{}, err
	}
	if config.HotShare, err = command.Flags().GetFloat64("hot-share"); err != nil {
		return Config{}, err
	}
	if config.StatsInterval, err = command.Flags().GetDuration("stats-interval"); err != nil {
		return Config{}, err
	}
	if config.Profile, err = command.Flags().GetString("profile"); err != nil {
		return Config{}, err
	}
	if config.OutputFormat, err = command.Flags().GetString("output-format"); err != nil {
		return Config{}, err
	}
	if config.DryRun, err = command.Flags().GetBool("dry-run"); err != nil {
		return Config{}, err
	}
	for _, item := range []struct {
		name   string
		target *time.Duration
	}{{"discovery-timeout", &config.DiscoveryTimeout}, {"mysql-dial-timeout", &config.DialTimeout}, {"mysql-read-timeout", &config.ReadTimeout}, {"mysql-write-timeout", &config.WriteTimeout}} {
		if *item.target, err = command.Flags().GetDuration(item.name); err != nil {
			return Config{}, err
		}
	}
	workloadConfig, err := command.Flags().GetString("workload-config")
	if err != nil {
		return Config{}, err
	}
	if workloadConfig != "" {
		if err := applyWorkloadFile(workloadConfig, &config, command.Flags().Changed); err != nil {
			return Config{}, err
		}
	}

	return config, config.Validate()
}
