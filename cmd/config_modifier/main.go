package main

import (
	"fmt"
	"log"
	"math"
	"metrics-bench-suite/pkg/samples"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ConfigModifier struct {
	ConfigPath    string
	OutputPath    string
	TargetSeries  int
	NumPods       int
	NumNodes      int
	NumInstances  int
	NumNamespaces int
	NumWorkloads  int
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
	"up",
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
	fileConfigs, err := samples.WalkAndParseConfig(c.ConfigPath)
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
			}
		}
		for _, tag := range fileConfig.Config.Tags {
			totalSeries *= len(tag.Dist.LabelGenerator().All())
		}

		maxNDV := math.Log(float64(c.TargetSeries)) / math.Log(float64(len(fileConfig.Config.Tags)))
		if slices.Contains(largeSeriesMetrics, fileConfig.Name) {
			maxNDV = 1
		}
		for i := range len(fileConfig.Config.Tags) {
			label := strings.ToLower(fileConfig.Config.Tags[i].Name)
			if _, exists := queryLabels[label]; !exists {
				fileConfig.Config.Tags[i].Dist.SetReplica(int(maxNDV), label)
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

func main() {
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
	rootCmd.Flags().IntP("num-pods", "p", 4, "The number of pods")
	rootCmd.Flags().IntP("num-nodes", "n", 8, "The number of nodes")
	rootCmd.Flags().IntP("num-instances", "i", 8, "The number of instances")
	rootCmd.Flags().IntP("num-namespaces", "N", 50, "The number of namespaces")
	rootCmd.Flags().IntP("num-workloads", "w", 4, "The number of workloads")
	rootCmd.Flags().IntP("target-series", "t", 10000, "The target time series number for each metric")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
