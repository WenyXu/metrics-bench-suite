package main

import (
	"fmt"
	"log"
	"metrics-bench-suite/pkg/timeseries"

	"github.com/spf13/cobra"
)

var logicalTbaleColumnTemplate = "`%s` STRING"
var logicalTableCreateTableSQLTemplate = `CREATE TABLE %s (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
%sTIME INDEX (greptime_timestamp),
PRIMARY KEY (%s)
) ENGINE = metric WITH (
on_physical_table = '%s'
);`

var physicalTbaleColumnTemplate = "`%s` STRING NULL SKIPPING INDEX WITH(granularity = '%d', type = 'BLOOM')"
var physicalTableCreateTableSQLTemplate = `CREATE TABLE %s (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
%sTIME INDEX (greptime_timestamp),
PRIMARY KEY (%s)
) ENGINE = metric WITH (
    "physical_metric_table" = "",   
    "memtable.type" = "partition_tree",
    "memtable.partition_tree.primary_key_encoding" = "sparse",
	"index.type" = "skipping", 
	"index.granularity" = "%d"
);`

// TableCreator is the struct for the table creator
type TableCreator struct {
	TableName                string
	TcpflowOutput            string
	SkippingIndexGranularity int
}

// Run is the entry point for the table creator
func (t *TableCreator) Run(cmd *cobra.Command, args []string) error {
	tableName, _ := cmd.Flags().GetString("table-name")
	t.TableName = tableName
	t.TcpflowOutput, _ = cmd.Flags().GetString("tcpflow-output")
	t.SkippingIndexGranularity, _ = cmd.Flags().GetInt("skipping-index-granularity")

	wrSet, err := timeseries.ParseAllRemoteWriteRequest(t.TcpflowOutput)
	if err != nil {
		return err
	}

	tsSet, err := timeseries.GetUniqueTimeSeries(wrSet)
	if err != nil {
		return err
	}

	labels := make(map[string]bool)
	for _, ts := range tsSet {
		for _, label := range ts.Labels {
			if label.Name == "__name__" {
				continue
			}
			if _, ok := labels[label.Name]; !ok {
				labels[label.Name] = true
			}
		}
	}

	var columnDef string
	var primaryKey string
	for label := range labels {
		columnDef += fmt.Sprintf(physicalTbaleColumnTemplate, label, t.SkippingIndexGranularity) + ",\n"
		primaryKey += fmt.Sprintf("`%s`,", label)
	}

	// create physical table
	fmt.Printf(physicalTableCreateTableSQLTemplate, t.TableName, columnDef, primaryKey, t.SkippingIndexGranularity)
	fmt.Printf("\n")

	createdLogicalTables := make(map[string]bool)
	// create logical tables
	for _, ts := range tsSet {
		var columnDef string
		var tableName string
		var primaryKey string
		for _, label := range ts.Labels {
			if label.Name == "__name__" {
				tableName = label.Value
				continue
			}
			columnDef += fmt.Sprintf(logicalTableCreateTableSQLTemplate, label.Name) + ",\n"
			primaryKey += fmt.Sprintf("`%s`,", label.Name)
		}

		if _, ok := createdLogicalTables[tableName]; !ok {
			fmt.Printf(logicalTableCreateTableSQLTemplate, fmt.Sprintf("`%s`", tableName), columnDef, primaryKey, t.TableName)
			fmt.Printf("\n")
			createdLogicalTables[tableName] = true
		}
	}

	return nil
}

func main() {
	tableCreator := &TableCreator{}

	var rootCmd = &cobra.Command{
		Use:   "table_creator",
		Short: "TableCreator is a tool to create table sql for the physical table and logical tables",
		Run: func(cmd *cobra.Command, args []string) {
			if err := tableCreator.Run(cmd, args); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}

	rootCmd.Flags().StringP("table-name", "t", "greptime_physical_table", "The name of the table")
	rootCmd.Flags().StringP("tcpflow-output", "o", "", "The path to the tcpflow output")
	rootCmd.Flags().IntP("skipping-index-granularity", "s", 102400, "The skipping index granularity")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
