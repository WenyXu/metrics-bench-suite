package config_modifier

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"metrics-bench-suite/pkg/samples"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ConfigModifier is a tool that modifies the config file to meet the target time series number.
type ConfigModifier struct {
	ConfigPath       string
	OutputPath       string
	TargetSeries     int
	NumPods          int
	NumNodes         int
	NumInstances     int
	NumNamespaces    int
	NumWorkloads     int
	ExpectTimeseries string
}

var queryLabels = map[string]string{
	"pod":           "pod",
	"node":          "node",
	"instance":      "instance",
	"namespace":     "namespace",
	"type":          "type",
	"cluster":       "cluster",
	"resource":      "resource",
	"job":           "job",
	"workload_type": "workload_type",
	"workload":      "workload",
	"device":        "device",
	"metrics_path":  "metrics_path",
	"le":            "le",
	"state":         "state",
	"code":          "code",
}

var largeSeriesMetrics = []string{
	"go_goroutines",
	"process_cpu_seconds_total",
	"process_resident_memory_bytes",
	"rest_client_requests_total",
	"process_resident_memory_bytes",
	"kube_node_status_allocatable",
}

// parseExpectTimeseries parses the expect timeseries file and returns a map of metric name to expected timeseries number
func parseExpectTimeseries(expectTimeseriesPath string) (map[string]int, error) {
	if expectTimeseriesPath == "" {
		return map[string]int{}, nil
	}

	file, err := os.Open(expectTimeseriesPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, "{}") {
			key := line[:len(line)-2]

			if !scanner.Scan() {
				return nil, fmt.Errorf("unexpected end of file")
			}

			valueStr := strings.TrimSpace(scanner.Text())
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing value: %s, %v", valueStr, err)
			}

			data[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return data, nil
}

func (c *ConfigModifier) run(cmd *cobra.Command, args []string) error {
	var err error
	c.ConfigPath, err = cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	c.OutputPath, err = cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	c.TargetSeries, err = cmd.Flags().GetInt("target-series")
	if err != nil {
		return err
	}
	c.NumPods, err = cmd.Flags().GetInt("num-pods")
	if err != nil {
		return err
	}
	c.NumNodes, err = cmd.Flags().GetInt("num-nodes")
	if err != nil {
		return err
	}
	c.NumInstances, err = cmd.Flags().GetInt("num-instances")
	if err != nil {
		return err
	}
	c.NumNamespaces, err = cmd.Flags().GetInt("num-namespaces")
	if err != nil {
		return err
	}
	c.NumWorkloads, err = cmd.Flags().GetInt("num-workloads")
	if err != nil {
		return err
	}
	c.ExpectTimeseries, err = cmd.Flags().GetString("expect-timeseries")
	if err != nil {
		return err
	}
	fileConfigs, err := samples.WalkAndParseConfig(c.ConfigPath)
	if err != nil {
		return err
	}

	expectTimeseries, err := parseExpectTimeseries(c.ExpectTimeseries)
	if err != nil {
		return err
	}

	allFilesTags := 0
	for _, fileConfig := range fileConfigs {
		totalSeries := 1

		for i := range len(fileConfig.Config.Tags) {
			if strings.ToLower(fileConfig.Config.Tags[i].Name) == "pod" {
				pod := "pod-"
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:          "replica_string",
					Replica:       &c.NumPods,
					ReplicaPrefix: &pod,
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "node" {
				node := "node-"
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:          "replica_string",
					Replica:       &c.NumNodes,
					ReplicaPrefix: &node,
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "instance" {
				instance := "instance-"
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:          "replica_string",
					Replica:       &c.NumInstances,
					ReplicaPrefix: &instance,
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "namespace" {
				namespace := "app-"
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:          "replica_string",
					Replica:       &c.NumNamespaces,
					ReplicaPrefix: &namespace,
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "workload" {
				workload := "workload-"
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:          "replica_string",
					Replica:       &c.NumWorkloads,
					ReplicaPrefix: &workload,
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "name" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "name0",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "id" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "id0",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "image" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "image0",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "pod_ip" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "10.0.0.0",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "node_ip" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "10.0.0.0",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "uid" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "f5a3af84-221f-486f-821f-1fbd61b5ede2",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "device" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "/dev/nvme0n1",
				}
			} else if strings.ToLower(fileConfig.Config.Tags[i].Name) == "container" {
				fileConfig.Config.Tags[i].Dist = samples.Distribution{
					Type:  "constant_string",
					Value: "container0",
				}
			}
		}
		for _, tag := range fileConfig.Config.Tags {
			totalSeries *= len(tag.Dist.LabelGenerator().All())
		}

		maxNDV := math.Log(float64(c.TargetSeries)) / math.Log(float64(len(fileConfig.Config.Tags)))
		if slices.Contains(largeSeriesMetrics, fileConfig.Name) {
			maxNDV = 1
		}
		if !strings.HasPrefix(fileConfig.Name, "ahead_") {
			for i := range len(fileConfig.Config.Tags) {
				label := strings.ToLower(fileConfig.Config.Tags[i].Name)
				if _, exists := queryLabels[label]; !exists {
					fileConfig.Config.Tags[i].Dist.SetReplica(int(maxNDV), label)
				}
			}
		}

		totalSeries = 1
		for _, tag := range fileConfig.Config.Tags {
			totalSeries *= len(tag.Dist.LabelGenerator().All())
		}
		allFilesTags += totalSeries
		if totalSeries > c.TargetSeries {
			log.Printf("Processing file: %s, total series: %d, tags: %d", fileConfig.Name, totalSeries, len(fileConfig.Config.Tags))
		}

		if expectedSeries, exists := expectTimeseries[fileConfig.Name]; exists {
			if totalSeries != expectedSeries {
				log.Printf("Processing file: %s, total series: %d, expected series: %d", fileConfig.Name, totalSeries, expectedSeries)
			}
		}
	}

	log.Printf("All files time series: %d", allFilesTags)

	err = os.MkdirAll(c.OutputPath, 0755)
	if err != nil {
		return err
	}

	for _, fileConfig := range fileConfigs {
		data, err := yaml.Marshal(fileConfig.Config)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%s/%s.yaml", c.OutputPath, fileConfig.Name)
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewCommand() *cobra.Command {
	configModifier := &ConfigModifier{}

	var rootCmd = &cobra.Command{
		Use:   "config_modifier",
		Short: "ConfigModifier is a tool to modify the config file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := configModifier.run(cmd, args); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
	rootCmd.Flags().StringP("config", "c", "", "The path to the config file")
	rootCmd.Flags().StringP("output", "o", "", "The output file")
	rootCmd.Flags().IntP("num-pods", "p", 2, "The number of pods")
	rootCmd.Flags().IntP("num-nodes", "n", 5, "The number of nodes")
	rootCmd.Flags().IntP("num-instances", "i", 5, "The number of instances")
	rootCmd.Flags().IntP("num-namespaces", "N", 50, "The number of namespaces")
	rootCmd.Flags().IntP("num-workloads", "w", 4, "The number of workloads")
	rootCmd.Flags().IntP("target-series", "t", 10000, "The target time series number for each metric")
	rootCmd.Flags().StringP("expect-timeseries", "e", "", "The expected time series number for each metric")

	return rootCmd
}
