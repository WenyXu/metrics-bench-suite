package partition_query_loader

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestBuildQueryPlansBuildsBoundedLogicalTableQuery(t *testing.T) {
	discovery := DiscoveryResult{Tables: []DiscoveredTable{{
		Database:      "metrics",
		LogicalTable:  "cpu_usage",
		PhysicalTable: "cpu_physical",
		ValueColumn:   "greptime_value",
		TimeIndex:     "greptime_timestamp",
		Partitions: []DiscoveredPartition{{
			Name:           "p0",
			RegionID:       42,
			LeaderDatanode: "datanode-1",
			Predicate:      Predicate{SQL: "`region` = ? AND `service` = ?", Args: []any{"us-east", "api"}},
		}},
	}}}

	result := BuildQueryPlans(discovery)
	if len(result.Plans) != 1 || len(result.Skipped) != 0 {
		t.Fatalf("unexpected plans result: %#v", result)
	}
	plan := result.Plans[0]
	wantSQL := "SELECT avg(`greptime_value`) FROM `metrics`.`cpu_usage` WHERE `region` = ? AND `service` = ? AND `greptime_timestamp` >= ? AND `greptime_timestamp` < ?"
	if plan.SQL != wantSQL {
		t.Fatalf("SQL = %q, want %q", plan.SQL, wantSQL)
	}
	if !strings.Contains(plan.SQL, "`region` = ?") || !strings.Contains(plan.SQL, "`greptime_timestamp` >= ?") || !strings.Contains(plan.SQL, "`greptime_timestamp` < ?") {
		t.Fatalf("query omitted a required bound: %q", plan.SQL)
	}
	start := time.Date(2026, time.July, 13, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Minute)
	if got, want := plan.Arguments(start, end), []any{"us-east", "api", start, end}; !reflect.DeepEqual(got, want) {
		t.Fatalf("arguments = %#v, want %#v", got, want)
	}
}

func TestBuildQueryPlansSkipsUnsafeTables(t *testing.T) {
	validPartition := DiscoveredPartition{Predicate: Predicate{SQL: "`region` = ?", Args: []any{"east"}}}
	tests := []struct {
		name   string
		table  DiscoveredTable
		reason string
	}{
		{name: "missing value column", table: DiscoveredTable{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", TimeIndex: "ts", Partitions: []DiscoveredPartition{validPartition}}, reason: "aggregate value"},
		{name: "missing time index", table: DiscoveredTable{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", ValueColumn: "value", Partitions: []DiscoveredPartition{validPartition}}, reason: "time index"},
		{name: "missing predicate", table: DiscoveredTable{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", ValueColumn: "value", TimeIndex: "ts"}, reason: "partition predicate"},
		{name: "mismatched predicate arguments", table: DiscoveredTable{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", ValueColumn: "value", TimeIndex: "ts", Partitions: []DiscoveredPartition{{Predicate: Predicate{SQL: "`region` = ?", Args: nil}}}}, reason: "partition predicate"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BuildQueryPlans(DiscoveryResult{Tables: []DiscoveredTable{test.table}})
			if len(result.Plans) != 0 || len(result.Skipped) != 1 {
				t.Fatalf("unexpected plans result: %#v", result)
			}
			if !strings.Contains(result.Skipped[0].Reason, test.reason) {
				t.Fatalf("skip reason = %q, want it to contain %q", result.Skipped[0].Reason, test.reason)
			}
		})
	}
}

func TestBuildQueryPlansQuotesDiscoveredIdentifiers(t *testing.T) {
	result := BuildQueryPlans(DiscoveryResult{Tables: []DiscoveredTable{{
		Database: "bad-name", LogicalTable: "cpu", PhysicalTable: "physical", ValueColumn: "value", TimeIndex: "ts",
		Partitions: []DiscoveredPartition{{Predicate: Predicate{SQL: "`region` = ?", Args: []any{"east"}}}},
	}}})
	if len(result.Plans) != 0 || len(result.Skipped) != 1 {
		t.Fatalf("unexpected plans result: %#v", result)
	}
}

func TestBuildQueryPlansQuotesColonSafeLogicalTable(t *testing.T) {
	result := BuildQueryPlans(DiscoveryResult{Tables: []DiscoveredTable{{
		Database: "metrics", LogicalTable: "cluster:cpu", PhysicalTable: "cpu_physical", ValueColumn: "value", TimeIndex: "ts",
		Partitions: []DiscoveredPartition{{Predicate: Predicate{SQL: "`region` = ?", Args: []any{"east"}}}},
	}}})
	if len(result.Plans) != 1 || !strings.Contains(result.Plans[0].SQL, "`metrics`.`cluster:cpu`") {
		t.Fatalf("colon table was not safely quoted: %#v", result)
	}
}
