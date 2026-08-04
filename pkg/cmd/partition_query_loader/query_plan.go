package partition_query_loader

import (
	"fmt"
	"time"
)

// QueryPlan is a parameterized, partition-targeted aggregate query on a logical table.
type QueryPlan struct {
	Database        string
	LogicalTable    string
	PhysicalTable   string
	Partition       string
	RegionID        uint64
	LeaderDatanode  string
	SQL             string
	MaxTimestampSQL string
	partitionArgs   []any
}

// Arguments returns the query arguments in SQL placeholder order for [start, end).
func (plan QueryPlan) Arguments(start, end time.Time) []any {
	args := make([]any, 0, len(plan.partitionArgs)+2)
	args = append(args, plan.partitionArgs...)
	return append(args, start, end)
}

// QueryPlanResult contains executable plans and the reasons unsafe tables were omitted.
type QueryPlanResult struct {
	Plans   []QueryPlan
	Skipped []SkipReason
}

// BuildQueryPlans converts discovered tables into bounded aggregate query templates.
func BuildQueryPlans(discovery DiscoveryResult) QueryPlanResult {
	result := QueryPlanResult{Skipped: append([]SkipReason(nil), discovery.Skipped...)}
	for _, table := range discovery.Tables {
		plans, reason := buildTableQueryPlans(table)
		if reason != "" {
			result.Skipped = append(result.Skipped, SkipReason{
				Database: table.Database,
				Table:    table.LogicalTable,
				Reason:   reason,
			})
			continue
		}
		result.Plans = append(result.Plans, plans...)
	}
	return result
}

func buildTableQueryPlans(table DiscoveredTable) ([]QueryPlan, string) {
	if !tableIdentifierPattern.MatchString(table.Database) || !tableIdentifierPattern.MatchString(table.LogicalTable) || !tableIdentifierPattern.MatchString(table.PhysicalTable) {
		return nil, "invalid discovered database or table identifier"
	}
	if !identifierPattern.MatchString(table.ValueColumn) {
		return nil, "no safe aggregate value column"
	}
	if !identifierPattern.MatchString(table.TimeIndex) {
		return nil, "no safe time index"
	}
	if len(table.Partitions) == 0 {
		return nil, "no partition predicate"
	}

	plans := make([]QueryPlan, 0, len(table.Partitions))
	for _, partition := range table.Partitions {
		predicate, err := safePartitionPredicate(partition.Predicate)
		if err != nil {
			return nil, fmt.Sprintf("no safe partition predicate: %v", err)
		}
		plans = append(plans, QueryPlan{
			Database:       table.Database,
			LogicalTable:   table.LogicalTable,
			PhysicalTable:  table.PhysicalTable,
			Partition:      partition.Name,
			RegionID:       partition.RegionID,
			LeaderDatanode: partition.LeaderDatanode,
			SQL: "SELECT avg(" + quoteIdentifier(table.ValueColumn) + ") FROM " +
				quoteIdentifier(table.Database) + "." + quoteIdentifier(table.LogicalTable) +
				" WHERE " + predicate.SQL + " AND " + quoteIdentifier(table.TimeIndex) +
				" >= ? AND " + quoteIdentifier(table.TimeIndex) + " < ?",
			MaxTimestampSQL: "SELECT " + quoteIdentifier(table.TimeIndex) + " FROM " + quoteIdentifier(table.Database) + "." + quoteIdentifier(table.LogicalTable) + " WHERE " + predicate.SQL + " ORDER BY " + quoteIdentifier(table.TimeIndex) + " DESC LIMIT 1",
			partitionArgs:   append([]any(nil), predicate.Args...),
		})
	}
	return plans, ""
}

func safePartitionPredicate(predicate Predicate) (Predicate, error) {
	if predicate.SQL == "" || len(predicate.Args) == 0 {
		return Predicate{}, fmt.Errorf("empty predicate")
	}
	if countPlaceholders(predicate.SQL) != len(predicate.Args) {
		return Predicate{}, fmt.Errorf("placeholder count does not match arguments")
	}
	return predicate, nil
}

func countPlaceholders(sql string) int {
	count := 0
	for _, char := range sql {
		if char == '?' {
			count++
		}
	}
	return count
}
