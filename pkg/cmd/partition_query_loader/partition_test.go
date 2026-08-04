package partition_query_loader

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParsePartitionDefinitionSupportsPlannedForms(t *testing.T) {
	tests := []struct {
		name       string
		definition string
		predicate  string
		args       []any
	}{
		{
			name: "single column equality",
			definition: "CREATE TABLE cpu (region STRING)\n" +
				"PARTITION ON COLUMNS (`region`) (region = 'region-1', region = 'region-2')",
			predicate: "`region` = ?",
			args:      []any{"region-1"},
		},
		{
			name:       "single column range",
			definition: "PARTITION ON COLUMNS (id) (id >= 100 AND id < 200)",
			predicate:  "`id` >= ? AND `id` < ?",
			args:       []any{"100", "200"},
		},
		{
			name:       "multi column conditions",
			definition: "PARTITION ON COLUMNS (`region`, service) (region = 'us-east' AND service = 'api')",
			predicate:  "`region` = ? AND `service` = ?",
			args:       []any{"us-east", "api"},
		},
		{
			name:       "grouped multi column conditions",
			definition: "PARTITION ON COLUMNS (`region`, service) ((region = 'us-east' AND service = 'api'))",
			predicate:  "`region` = ? AND `service` = ?",
			args:       []any{"us-east", "api"},
		},
		{
			name:       "timestamp literal",
			definition: "PARTITION ON COLUMNS (ts) (ts >= TIMESTAMP '2025-01-02 03:04:05')",
			predicate:  "`ts` >= ?",
			args:       []any{time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			definition, err := ParsePartitionDefinition(test.definition)
			if err != nil {
				t.Fatalf("parse partition definition: %v", err)
			}
			predicate, err := BuildPartitionPredicate(definition.Partitions[0])
			if err != nil {
				t.Fatalf("build predicate: %v", err)
			}
			if predicate.SQL != test.predicate {
				t.Fatalf("expected predicate %q, got %q", test.predicate, predicate.SQL)
			}
			if !reflect.DeepEqual(predicate.Args, test.args) {
				t.Fatalf("expected arguments %#v, got %#v", test.args, predicate.Args)
			}
		})
	}
}

func TestParsePartitionDefinitionPreservesExactNumericLiteral(t *testing.T) {
	definition, err := ParsePartitionDefinition("PARTITION ON COLUMNS (id) (id = 9007199254740993.000000000000000001)")
	if err != nil {
		t.Fatal(err)
	}
	if value := definition.Partitions[0].Conditions[0].Value; value != "9007199254740993.000000000000000001" {
		t.Fatalf("numeric value was changed: %#v", value)
	}
}

func TestBuildPartitionPredicateNeverInterpolatesLiteral(t *testing.T) {
	definition, err := ParsePartitionDefinition("PARTITION ON COLUMNS (region) (region = 'x'' OR 1=1 --')")
	if err != nil {
		t.Fatalf("parse partition definition: %v", err)
	}

	predicate, err := BuildPartitionPredicate(definition.Partitions[0])
	if err != nil {
		t.Fatalf("build predicate: %v", err)
	}
	if predicate.SQL != "`region` = ?" {
		t.Fatalf("expected parameterized predicate, got %q", predicate.SQL)
	}
	if strings.Contains(predicate.SQL, "OR 1=1") {
		t.Fatalf("predicate interpolated a literal: %q", predicate.SQL)
	}
	if !reflect.DeepEqual(predicate.Args, []any{"x' OR 1=1 --"}) {
		t.Fatalf("unexpected predicate arguments: %#v", predicate.Args)
	}
}

func TestParsePartitionDefinitionRejectsUnsafeOrUnsupportedInput(t *testing.T) {
	tests := []string{
		"CREATE TABLE no_partition (id INT)",
		"PARTITION ON COLUMNS (id) (id = ?)",
		"PARTITION ON COLUMNS (id) (id IN (1, 2))",
		"PARTITION ON COLUMNS (region, service) (region = 'us-east')",
		"PARTITION ON COLUMNS (id) (unknown = 1)",
		"PARTITION ON COLUMNS (id) ()",
	}

	for _, definitionText := range tests {
		if _, err := ParsePartitionDefinition(definitionText); err == nil {
			t.Fatalf("expected unsupported definition to be skipped: %q", definitionText)
		}
	}
}

func TestBuildPartitionPredicateRejectsInvalidCondition(t *testing.T) {
	_, err := BuildPartitionPredicate(Partition{Conditions: []PartitionCondition{{
		Column:   "region; DROP TABLE metrics",
		Operator: "=",
		Value:    "us-east",
	}}})
	if err == nil {
		t.Fatal("expected invalid partition condition to be rejected")
	}
}
