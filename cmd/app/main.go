package main

import (
	"log"
	"metrics-bench-suite/pkg/cmd/config_modifier"
	"metrics-bench-suite/pkg/cmd/loader"
	"metrics-bench-suite/pkg/cmd/remote_write_request_viewer"
	"metrics-bench-suite/pkg/cmd/sample_generator"
	"metrics-bench-suite/pkg/cmd/sample_loader"
	"metrics-bench-suite/pkg/cmd/schema_generator"
	"metrics-bench-suite/pkg/cmd/table_creator"
	"metrics-bench-suite/pkg/cmd/timeseries_analyzer"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var sampleLoaderCmd = sample_loader.NewCommand()
	var sampleGeneratorCmd = sample_generator.NewCommand()
	var configModifierCmd = config_modifier.NewCommand()
	var remoteWriteRequestViewerCmd = remote_write_request_viewer.NewCommand()
	var loaderCmd = loader.NewCommand()
	var schemaGeneratorCmd = schema_generator.NewCommand()
	var tableCreatorCmd = table_creator.NewCommand()
	var timeseriesAnalyzerCmd = timeseries_analyzer.NewCommand()
	rootCmd.AddCommand(configModifierCmd)
	rootCmd.AddCommand(sampleLoaderCmd)
	rootCmd.AddCommand(sampleGeneratorCmd)
	rootCmd.AddCommand(remoteWriteRequestViewerCmd)
	rootCmd.AddCommand(loaderCmd)
	rootCmd.AddCommand(schemaGeneratorCmd)
	rootCmd.AddCommand(tableCreatorCmd)
	rootCmd.AddCommand(timeseriesAnalyzerCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
