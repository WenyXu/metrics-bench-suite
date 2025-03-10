package main

import (
	"fmt"
	"log"
	"metrics-bench-suite/pkg/samples"
	"os"

	"github.com/spf13/cobra"
)

var logicalTbaleColumnTemplate = "`%s` STRING"
var logicalTableCreateTableSQLTemplate = `
CREATE TABLE %s (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
%sTIME INDEX (greptime_timestamp),
PRIMARY KEY (%s)
) ENGINE = metric WITH (
on_physical_table = '%s'
);
`

var logicalTableCreateTableSQLTemplateWithoutPrimaryKey = `
CREATE TABLE %s (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
%sTIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = '%s'
);
`

var physicalTbaleColumnTemplate = "`%s` STRING NULL SKIPPING INDEX WITH(granularity = '%d', type = 'BLOOM')"
var physicalTableCreateTableSQLTemplate = `
CREATE TABLE %s (
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
);
`

// TableCreator is the struct for the table creator
type TableCreator struct {
	TableName                string
	ConfigPath               string
	SkippingIndexGranularity int
	Database                 string
}

// Run is the entry point for the table creator
func (t *TableCreator) Run(cmd *cobra.Command, args []string) error {
	tableName, _ := cmd.Flags().GetString("table-name")
	t.TableName = tableName
	t.ConfigPath, _ = cmd.Flags().GetString("config")
	t.SkippingIndexGranularity, _ = cmd.Flags().GetInt("skipping-index-granularity")
	t.Database, _ = cmd.Flags().GetString("database")
	fileConfigs, err := samples.WalkAndParseConfig(t.ConfigPath)
	if err != nil {
		return err
	}
	if len(fileConfigs) == 0 {
		return fmt.Errorf("no config files found")
	}

	labels := make(map[string]bool)
	for _, fileConfig := range fileConfigs {
		for _, tag := range fileConfig.Config.Tags {
			if _, ok := labels[tag.Name]; !ok {
				labels[tag.Name] = true
			}
		}
	}

	fileName := fmt.Sprintf("%s.%s-create-tables.sql", t.Database, t.TableName)
	log.Printf("Writing to file: %s", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;\n", t.Database)); err != nil {
		return err
	}

	var columnDef string
	var primaryKey string
	for label := range labels {
		columnDef += fmt.Sprintf(physicalTbaleColumnTemplate, label, t.SkippingIndexGranularity) + ",\n"
		primaryKey += fmt.Sprintf("`%s`,", label)
	}

	// create physical table
	if _, err := file.WriteString(fmt.Sprintf(physicalTableCreateTableSQLTemplate, fmt.Sprintf("`%s`.`%s`", t.Database, t.TableName), columnDef, primaryKey, t.SkippingIndexGranularity)); err != nil {
		return err
	}

	createdLogicalTables := make(map[string]bool)
	// create logical tables
	for _, fileConfig := range fileConfigs {
		var columnDef string
		var primaryKey string
		for _, label := range fileConfig.Config.Tags {
			columnDef += fmt.Sprintf(logicalTbaleColumnTemplate, label.Name) + ",\n"
			primaryKey += fmt.Sprintf("`%s`,", label.Name)
		}

		if _, ok := createdLogicalTables[fileConfig.Name]; !ok {
			if len(fileConfig.Config.Tags) == 0 {
				if _, err := file.WriteString(fmt.Sprintf(logicalTableCreateTableSQLTemplateWithoutPrimaryKey, fmt.Sprintf("`%s`.`%s`", t.Database, fileConfig.Name), columnDef, t.TableName)); err != nil {
					return err
				}
			} else {
				if _, err := file.WriteString(fmt.Sprintf(logicalTableCreateTableSQLTemplate, fmt.Sprintf("`%s`.`%s`", t.Database, fileConfig.Name), columnDef, primaryKey, t.TableName)); err != nil {
					return err
				}
			}
			createdLogicalTables[fileConfig.Name] = true
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

	rootCmd.Flags().StringP("database", "d", "public", "The name of the database")
	rootCmd.Flags().StringP("table-name", "t", "greptime_physical_table", "The name of the table")
	rootCmd.Flags().StringP("config", "c", "", "The path to the config files")
	rootCmd.Flags().IntP("skipping-index-granularity", "s", 102400, "The skipping index granularity")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
