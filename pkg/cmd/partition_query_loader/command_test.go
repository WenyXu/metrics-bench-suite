package partition_query_loader

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"
)

func TestRunConfiguredDryRunOrchestratesDiscoveryAndPlans(t *testing.T) {
	config := DefaultConfig()
	config.DryRun = true
	var discovered, ran, sampled, wrote bool
	dependencies := workflowDependencies{
		openDB: func(Config) (*sql.DB, error) { return nil, nil },
		discover: func(context.Context, *sql.DB, []string) (DiscoveryResult, error) {
			discovered = true
			return DiscoveryResult{Tables: []DiscoveredTable{{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", ValueColumn: "value", TimeIndex: "ts", Partitions: []DiscoveredPartition{{Name: "p0", RegionID: 7, Predicate: Predicate{SQL: "`region` = ?", Args: []any{"east"}}}}}}}, nil
		},
		runWorkload: func(context.Context, *sql.DB, Config, []QueryPlan) (WorkloadResult, error) {
			ran = true
			return WorkloadResult{}, nil
		},
		sample: func(context.Context, *sql.DB, []DiscoveredTable, time.Time) ([]RegionSnapshot, error) {
			sampled = true
			return nil, nil
		},
		write: func(io.Writer, string, []Observation) error { wrote = true; return nil },
		now:   time.Now,
	}
	var output bytes.Buffer
	result, err := runConfigured(context.Background(), config, &output, dependencies)
	if err != nil {
		t.Fatal(err)
	}
	if !discovered || ran || sampled || wrote {
		t.Fatalf("unexpected workflow calls: discovery=%t workload=%t sample=%t write=%t", discovered, ran, sampled, wrote)
	}
	if len(result.ByRegion) != 0 || output.String() != "region_id=7 partition=p0 sql=SELECT avg(`value`) FROM `metrics`.`cpu` WHERE `region` = ? AND `ts` >= ? AND `ts` < ?\n" {
		t.Fatalf("result=%#v output=%q", result, output.String())
	}
}

func TestNewCommandUsesCommandContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	command := NewCommand()
	command.SetContext(ctx)
	command.SetArgs([]string{"--mysql-user", "user"})
	if err := command.Execute(); err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected command context cancellation, got %v", err)
	}
}
