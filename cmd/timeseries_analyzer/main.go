package main

import (
	"log"
	"metrics-bench-suite/pkg/parser"
	"metrics-bench-suite/pkg/timeseries"
	"metrics-bench-suite/pkg/utils/decode"

	"github.com/spf13/cobra"
)

// Analyzer is the struct for the timeseries analyzer
type Analyzer struct{}

// Run is the main function to run the timeseries analyzer
func (a *Analyzer) Run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		log.Fatal("You must provide input data")
	}

	filePath := args[0]
	log.Printf("Processing file: %s", filePath)

	tableNameCounter := make(map[string]bool)
	timeSeriesCounter := make(map[string]bool)

	requests, err := parser.ParseHTTPRequests(filePath)
	if err != nil {
		log.Printf("Error parsing request file %s: %v", filePath, err)
		return err
	}

	for _, request := range requests {
		wr, err := decode.Body(request.Body)
		if err != nil {
			log.Fatalf("failed to decode body: %v", err)
		}
		err = timeseries.CountTimeSeries(&wr, &tableNameCounter, &timeSeriesCounter)
		if err != nil {
			log.Fatalf("failed to count time series: %v", err)
		}
	}

	log.Printf("Found %d table names", len(tableNameCounter))
	log.Printf("Found %d time series", len(timeSeriesCounter))

	return nil
}

func main() {
	analyzer := &Analyzer{}

	var rootCmd = &cobra.Command{
		Use:   "timeseries_analyzer <tcpflow_output_dir>",
		Short: "Timeseries Analyzer is a tool to analyze timeseries in tcpflow output",
		Run: func(cmd *cobra.Command, args []string) {
			analyzer.Run(cmd, args)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error processing tcpflow files: %v", err)
	}
}
