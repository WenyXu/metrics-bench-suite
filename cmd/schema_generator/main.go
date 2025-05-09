package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Generate an array of arrays, where each inner array contains m unique random numbers
func generateRandomArrays(n int, mRange [2]int, count int) [][]int {
	rand.NewSource(time.Now().UnixNano())
	arrays := make([][]int, count)
	for i := range count {
		m := rand.Intn(mRange[1]-mRange[0]+1) + mRange[0]
		arr := rand.Perm(n)[:m]
		arrays[i] = arr
	}
	return arrays
}

// Generate SQL for creating a metrics table
func generateMetricsSQL(columnListCount, sampleCount int, columnCountRange [2]int, physicalTables []string) ([]map[string]interface{}, string) {
	sqlTemplate := "CREATE TABLE IF NOT EXISTS `%s` (\n`greptime_timestamp` TIMESTAMP(3) NOT NULL,\n`greptime_value` DOUBLE NULL,\n`region` STRING NULL,\n%s,\nTIME INDEX (greptime_timestamp),\nPRIMARY KEY (`region`,%s),\n) ENGINE = metric WITH (\non_physical_table = '%s'\n);"

	var tableInfoList []map[string]any
	var finalSQLList []string
	currentTableID := 0

	for _, physicalTable := range physicalTables {
		selectedColumns := generateRandomArrays(columnListCount, columnCountRange, sampleCount)
		for _, columns := range selectedColumns {
			tableName := fmt.Sprintf("metrics_table_%d", currentTableID)
			var columnsStr, primaryKeys []string
			for _, col := range columns {
				columnsStr = append(columnsStr, fmt.Sprintf("`column%d` STRING NULL", col))
				primaryKeys = append(primaryKeys, fmt.Sprintf("`column%d`", col))
			}
			tableInfo := map[string]any{
				"table_name":     tableName,
				"columns":        columns,
				"physical_table": physicalTable,
			}
			tableInfoList = append(tableInfoList, tableInfo)
			sqlStr := fmt.Sprintf(sqlTemplate, tableName, strings.Join(columnsStr, ",\n"), strings.Join(primaryKeys, ","), physicalTable)
			finalSQLList = append(finalSQLList, sqlStr)
			currentTableID++
		}
	}
	return tableInfoList, strings.Join(finalSQLList, "\n")
}

// Generate SQL for creating a physical table
func generatePhysicalTableSQL(count int) ([]string, string) {
	sqlTemplate := "CREATE TABLE IF NOT EXISTS `%s` (\n  `greptime_timestamp` TIMESTAMP(3) NOT NULL,\n    `greptime_value` DOUBLE NULL,\n    `region` STRING NULL,\n    `column1` STRING NULL,\n    `column2` STRING NULL,\n    `column3` STRING NULL,\n    `column4` STRING NULL,\n    `column5` STRING NULL,\n    `column6` STRING NULL,\n    `column7` STRING NULL,\n    `column8` STRING NULL,\n    `column9` STRING NULL,\n    `column10` STRING NULL,\n    `column11` STRING NULL,\n    `column12` STRING NULL,\n    `column13` STRING NULL,\n    `column14` STRING NULL,\n    `column15` STRING NULL,\n    `column16` STRING NULL,\n    `column17` STRING NULL,\n    `column18` STRING NULL,\n    `column19` STRING NULL,\n    `column20` STRING NULL,\n    `column21` STRING NULL,\n    `column22` STRING NULL,\n    `column23` STRING NULL,\n    `column24` STRING NULL,\n    `column25` STRING NULL,\n    `column26` STRING NULL,\n    `column27` STRING NULL,\n    `column28` STRING NULL,\n    `column29` STRING NULL,\n    `column30` STRING NULL,\n    `column31` STRING NULL,\n    `column32` STRING NULL,\n    `column33` STRING NULL,\n    `column34` STRING NULL,\n    `column35` STRING NULL,\n    `column36` STRING NULL,\n    `column37` STRING NULL,\n    `column38` STRING NULL,\n    `column39` STRING NULL,\n    `column40` STRING NULL,\n    `column41` STRING NULL,\n    `column42` STRING NULL,\n    `column43` STRING NULL,\n    `column44` STRING NULL,\n    `column45` STRING NULL,\n    `column46` STRING NULL,\n    `column47` STRING NULL,\n    `column48` STRING NULL,\n    `column49` STRING NULL,\n    `column50` STRING NULL,\n    `column51` STRING NULL,\n    `column52` STRING NULL,\n    `column53` STRING NULL,\n    `column54` STRING NULL,\n    `column55` STRING NULL,\n    `column56` STRING NULL,\n    `column57` STRING NULL,\n    `column58` STRING NULL,\n    `column59` STRING NULL,\n    `column60` STRING NULL,\n    `column61` STRING NULL,\n    `column62` STRING NULL,\n    `column63` STRING NULL,\n    `column64` STRING NULL,\n    `column65` STRING NULL,\n    `column66` STRING NULL,\n    `column67` STRING NULL,\n    `column68` STRING NULL,\n    `column69` STRING NULL,\n    `column70` STRING NULL,\n    `column71` STRING NULL,\n    `column72` STRING NULL,\n    `column73` STRING NULL,\n    `column74` STRING NULL,\n    `column75` STRING NULL,\n    `column76` STRING NULL,\n    `column77` STRING NULL,\n    `column78` STRING NULL,\n    `column79` STRING NULL,\n    `column80` STRING NULL,\n    `column81` STRING NULL,\n    `column82` STRING NULL,\n    `column83` STRING NULL,\n    `column84` STRING NULL,\n    `column85` STRING NULL,\n    `column86` STRING NULL,\n    `column87` STRING NULL,\n    `column88` STRING NULL,\n    `column89` STRING NULL,\n    `column90` STRING NULL,\n    `column91` STRING NULL,\n    `column92` STRING NULL,\n    `column93` STRING NULL,\n    `column94` STRING NULL,\n    `column95` STRING NULL,\n    `column96` STRING NULL,\n    `column97` STRING NULL,\n    `column98` STRING NULL,\n    `column99` STRING NULL,\n    TIME INDEX (`greptime_timestamp`),\n    PRIMARY KEY (`region`,`column1`, `column2`, `column3`, `column4`, `column5`, `column6`, `column7`, `column8`, `column9`, `column10`,\n      `column11`, `column12`, `column13`, `column14`, `column15`, `column16`, `column17`, `column18`, `column19`, `column20`,\n      `column21`, `column22`, `column23`, `column24`, `column25`, `column26`, `column27`, `column28`, `column29`, `column30`,\n      `column31`, `column32`, `column33`, `column34`, `column35`, `column36`, `column37`, `column38`, `column39`, `column40`,\n      `column41`, `column42`, `column43`, `column44`, `column45`, `column46`, `column47`, `column48`, `column49`, `column50`,\n      `column51`, `column52`, `column53`, `column54`, `column55`, `column56`, `column57`, `column58`, `column59`, `column60`,\n      `column61`, `column62`, `column63`, `column64`, `column65`, `column66`, `column67`, `column68`, `column69`, `column70`,\n      `column71`, `column72`, `column73`, `column74`, `column75`, `column76`, `column77`, `column78`, `column79`, `column80`,\n      `column81`, `column82`, `column83`, `column84`, `column85`, `column86`, `column87`, `column88`, `column89`, `column90`,\n      `column91`, `column92`, `column93`, `column94`, `column95`, `column96`, `column97`, `column98`, `column99`)\n)\nPARTITION ON COLUMNS (region) (\n  region = 'region-0',\n  region = 'region-1',\n  region = 'region-2',\n  region = 'region-3',\n  region = 'region-4',\n  region = 'region-5',\n  region = 'region-6',\n  region = 'region-7',\n  region = 'region-8',\n  region = 'region-9',\n  region = 'region-10',\n  region = 'region-11',\n  region = 'region-12',\n  region = 'region-13',\n  region = 'region-14',\n  region = 'region-15',\n  region = 'region-16',\n  region = 'region-17',\n  region = 'region-18',\n  region = 'region-19',\n  region = 'region-20',\n  region = 'region-21',\n  region = 'region-22',\n  region = 'region-23',\n  region = 'region-24',\n  region = 'region-25',\n  region = 'region-26',\n  region = 'region-27',\n  region = 'region-28',\n  region = 'region-29',\n  region = 'region-30',\n  region = 'region-31',\n  region = 'region-32',\n  region = 'region-33',\n  region = 'region-34',\n  region = 'region-35',\n  region = 'region-36',\n  region = 'region-37',\n  region = 'region-38',\n  region = 'region-39',\n  region = 'region-40',\n  region = 'region-41',\n  region = 'region-42',\n  region = 'region-43',\n  region = 'region-44',\n  region = 'region-45',\n  region = 'region-46',\n  region = 'region-47',\n  region = 'region-48',\n  region = 'region-49',\n  region = 'region-50',\n)\nENGINE=metric\nWITH(\n  physical_metric_table = 'true'\n);"

	var sqlList []string
	var tableNameList []string

	for i := range count {
		tableName := fmt.Sprintf("table_%d", i)
		sqlStr := fmt.Sprintf(sqlTemplate, tableName)
		sqlList = append(sqlList, sqlStr)
		tableNameList = append(tableNameList, tableName)
	}
	return tableNameList, strings.Join(sqlList, "\n")
}

// Write content to a file
func writeToFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

func generateSampleLoaderYaml(tableInfos *[]map[string]any, targetFilePath string) error {
	// create a path to store yaml files
	err := os.MkdirAll(targetFilePath, os.ModePerm)
	if err != nil {
		return err
	}
	for _, tableInfo := range *tableInfos {
		tableName := tableInfo["table_name"].(string)
		columns := tableInfo["columns"].([]int)
		yamlFileName := fmt.Sprintf("%s/%s.yaml", targetFilePath, tableName)
		tags := make([]map[string]any, len(columns)+1)
		for i, col := range columns {
			tags[i] = map[string]any{
				"name": fmt.Sprintf("column%d", col),
				"type": "string",
				"dist": map[string]any{
					"type":  "constant_string",
					"value": fmt.Sprintf("value%d", col),
				},
			}
		}
		tags[len(columns)] = map[string]any{
			"name": "region",
			"type": "string",
			"dist": map[string]any{
				"type":           "replica_string",
				"replica":        50,
				"replica_prefix": "region-",
			},
		}

		yamlData := map[string]interface{}{
			"start":    "2023-01-01T00:00:00Z",
			"end":      "2023-01-02T00:00:00Z",
			"interval": 30,
			"tags":     tags,
			"fields": []map[string]any{
				{
					"name": "greptime_value",
					"type": "float",
					"dist": map[string]any{
						"type":        "mono_inc",
						"lower_bound": 1,
						"upper_bound": 100,
						"step":        10,
					},
				},
			},
		}
		yamlContent, err := yaml.Marshal(yamlData)
		if err != nil {
			return err
		}
		err = writeToFile(yamlFileName, string(yamlContent))
		if err != nil {
			return err
		}
	}
	return nil
}

func run(columnListCount, sampleCount int, columnCountRange [2]int, targetPath string) error {

	physicalTableSQLFileName := fmt.Sprintf("%s/physical_table.sql", targetPath)
	metricsTableSQLFileName := fmt.Sprintf("%s.metrics_table.sql", targetPath)

	// Generate SQL for creating a physical table
	physicalTableNameList, physicalTableSQL := generatePhysicalTableSQL(400)
	err := writeToFile(physicalTableSQLFileName, physicalTableSQL)
	if err != nil {
		fmt.Printf("Error writing physical table SQL: %v\n", err)
		return err
	}

	// Generate SQL for creating a metrics table
	tableInfos, metricsTableSQL := generateMetricsSQL(columnListCount, sampleCount, columnCountRange, physicalTableNameList)
	err = writeToFile(metricsTableSQLFileName, metricsTableSQL)
	if err != nil {
		fmt.Printf("Error writing metrics table SQL: %v\n", err)
		return err
	}

	// Save table info to a file
	tableInfoFile, err := os.Create(fmt.Sprintf("%s/metrics_table_info.json", targetPath))
	if err != nil {
		fmt.Printf("Error creating table info file: %v\n", err)
		return err
	}
	defer tableInfoFile.Close()
	jsonEncoder := json.NewEncoder(tableInfoFile)
	jsonEncoder.SetIndent("", "  ")
	err = jsonEncoder.Encode(tableInfos)
	if err != nil {
		fmt.Printf("Error writing table info: %v\n", err)
		return err
	}
	err = generateSampleLoaderYaml(&tableInfos, targetPath)
	if err != nil {
		fmt.Printf("Error generating sample loader YAML: %v\n", err)
		return err
	}
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "schema_generator",
		Short: "SchemaGenerator is a tool to generate SQL schema",
		Run: func(cmd *cobra.Command, args []string) {
			columnListCount, err := cmd.Flags().GetInt("column-list-count")
			if err != nil {
				log.Fatalf("Error getting column-list-count: %v", err)
			}
			selectedColumnCount, err := cmd.Flags().GetString("selected-column-count-range")
			if err != nil {
				log.Fatalf("Error getting selected-column-count-range: %v", err)
			}
			selectedColumnCountRange := strings.Split(selectedColumnCount, ",")
			if len(selectedColumnCountRange) != 2 {
				log.Fatalf("Invalid selected-column-count-range range: %s", selectedColumnCount)
			}
			columnCountRange := [2]int{}
			columnCountRange[0], err = strconv.Atoi(selectedColumnCountRange[0])
			if err != nil {
				log.Fatalf("Error converting selected-column-count-range range: %v", err)
			}
			columnCountRange[1], err = strconv.Atoi(selectedColumnCountRange[1])
			if err != nil {
				log.Fatalf("Error converting selected-column-count-range range: %v", err)
			}
			if columnCountRange[0] > columnCountRange[1] {
				log.Fatalf("Invalid selected-column-count-range range: %s", selectedColumnCount)
			}
			sampleCount, err := cmd.Flags().GetInt("sample-count")
			if err != nil {
				log.Fatalf("Error getting sample-count: %v", err)
			}
			targetPath, err := cmd.Flags().GetString("target-path")
			if err != nil {
				log.Fatalf("Error getting target-path: %v", err)
			}
			log.Printf("Generating schema with column list count: %d, sample count: %d, column count range: %v, target path: %s", columnListCount, sampleCount, columnCountRange, targetPath)
			err = run(columnListCount, sampleCount, columnCountRange, targetPath)
			if err != nil {
				log.Fatalf("Error running sample_loader: %v", err)
			}
		},
	}
	rootCmd.Flags().IntP("column-list-count", "c", 100, "The number of columns in the column list")
	rootCmd.Flags().StringP("selected-column-count-range", "s", "5,10", "The range of selected column count")
	rootCmd.Flags().IntP("sample-count", "n", 10, "The number of samples to generate")
	rootCmd.Flags().StringP("target-path", "t", "yaml_config", "The target path to save the generated files")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
