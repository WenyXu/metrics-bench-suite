package loader

import (
	"fmt"
	"log"
	"metrics-bench-suite/pkg/http"
	"metrics-bench-suite/pkg/parser"
	"metrics-bench-suite/pkg/timeseries"
	"metrics-bench-suite/pkg/utils/decode"
	"sync"
	"time"

	"math/rand"
	"slices"

	"github.com/prometheus/prometheus/prompb"
	"github.com/spf13/cobra"
)

// Loader is the struct for the loader
type Loader struct {
	URL                  string
	TcpflowOutput        string
	TimeseriesPerRequest int
	Interval             time.Duration
	MockInterval         time.Duration
	StartDate            time.Time
	DryRun               bool
	SampleScale          int
	MetricScale          int
	LabelScale           int
	Seed                 int
}

// Run is the main function for the loader
func (l *Loader) Run(cmd *cobra.Command, args []string) error {
	intervalStr, _ := cmd.Flags().GetString("interval")
	mockIntervalStr, _ := cmd.Flags().GetString("mock-interval")
	initialDateStr, _ := cmd.Flags().GetString("start-date")
	l.URL, _ = cmd.Flags().GetString("url")
	l.TcpflowOutput, _ = cmd.Flags().GetString("tcpflow-output")
	l.TimeseriesPerRequest, _ = cmd.Flags().GetInt("timeseries-per-request")
	l.Seed, _ = cmd.Flags().GetInt("seed")
	l.DryRun, _ = cmd.Flags().GetBool("dry-run")
	l.SampleScale, _ = cmd.Flags().GetInt("sample-scale")
	l.MetricScale, _ = cmd.Flags().GetInt("metric-scale")
	l.LabelScale, _ = cmd.Flags().GetInt("label-scale")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}
	l.MockInterval, err = time.ParseDuration(mockIntervalStr)
	if err != nil {
		return err
	}
	initialDate, err := time.Parse(time.RFC3339, initialDateStr)
	if err != nil {
		return err
	}
	l.StartDate = initialDate
	l.Interval = interval
	log.Printf("Start date: %s", l.StartDate)
	log.Printf("Interval: %s", l.Interval)
	log.Printf("Mock interval: %s", l.MockInterval)
	log.Printf("Seed: %d", l.Seed)
	log.Printf("Timeseries per request: %d", l.TimeseriesPerRequest)
	log.Printf("Sample scale: %d", l.SampleScale)
	log.Printf("Metric scale: %d", l.MetricScale)
	log.Printf("Label scale: %d", l.LabelScale)
	r := rand.New(rand.NewSource(int64(l.Seed)))

	wrSet, err := l.getAllRemoteWriteRequest()
	if err != nil {
		return err
	}

	err = l.printStats(&wrSet)
	if err != nil {
		return err
	}

	tsSet, err := l.getTimeSeries(wrSet)
	if err != nil {
		return err
	}

	fs, err := l.buildFactorySet(tsSet, r)
	if err != nil {
		return err
	}

	if l.DryRun {
		log.Printf("Dry run, skipping processing")
		return nil
	}

	ticker := time.NewTicker(l.Interval)
	defer ticker.Stop()

	for range ticker.C {
		tsSet, err := fs.Generate(l.SampleScale)
		if err != nil {
			return err
		}
		tsSet = ScaleMetrics(tsSet, l.MetricScale)
		tsSet = ScaleLabels(tsSet, l.LabelScale)
		err = l.process(tsSet)
		if err != nil {
			log.Printf("failed to send write request: %v", err)
		}
	}

	return nil
}

// CloneTimeSeries clones a time series
func CloneTimeSeries(ts prompb.TimeSeries) prompb.TimeSeries {
	clone := prompb.TimeSeries{}
	clone.Labels = make([]prompb.Label, len(ts.Labels))
	copy(clone.Labels, ts.Labels)
	clone.Samples = make([]prompb.Sample, len(ts.Samples))
	copy(clone.Samples, ts.Samples)
	return clone
}

// ScaleMetrics scales the metrics by the given scale
func ScaleMetrics(tsSet []prompb.TimeSeries, scale int) []prompb.TimeSeries {
	total := len(tsSet)
	for i := 1; i < scale; i++ {
		for _, ts := range tsSet[:total] {
			cloneTs := CloneTimeSeries(ts)
			cloneTs.Labels[0].Value = fmt.Sprintf("%s_%d", ts.Labels[0].Value, i)
			tsSet = append(tsSet, cloneTs)
		}
	}

	return tsSet
}

// ScaleLabels scales the labels by the given scale
func ScaleLabels(tsSet []prompb.TimeSeries, scale int) []prompb.TimeSeries {
	total := len(tsSet)
	for i := 1; i < scale; i++ {
		for _, ts := range tsSet[:total] {
			cloneTs := CloneTimeSeries(ts)
			for j := 0; j < len(cloneTs.Labels); j++ {
				if cloneTs.Labels[j].Name == "__name__" {
					continue
				}
				cloneTs.Labels[j].Name = fmt.Sprintf("%s%c", cloneTs.Labels[j].Name, 96+i)
			}
			tsSet = append(tsSet, cloneTs)
		}
	}

	return tsSet
}

func (l *Loader) sendWriteRequest(tsSet []prompb.TimeSeries) error {
	r := http.NewRequester(l.URL)
	wr := prompb.WriteRequest{
		Timeseries: tsSet,
	}
	return r.Send(wr)
}

func (l *Loader) process(tsSet []prompb.TimeSeries) error {
	wg := sync.WaitGroup{}
	now := time.Now()
	total := len(tsSet)
	chunks := slices.Collect(slices.Chunk(tsSet, l.TimeseriesPerRequest))
	log.Printf("Sending %d chunks, chunk size: %d, total time series: %d, total samples: %d", len(chunks), l.TimeseriesPerRequest, total/l.SampleScale, total)
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []prompb.TimeSeries) {
			defer wg.Done()
			err := l.sendWriteRequest(chunk)
			if err != nil {
				log.Printf("failed to send write request: %v", err)
			}
		}(chunk)
	}
	wg.Wait()
	log.Printf("Processed %d time series in %s", total, time.Since(now))
	return nil
}

func (l *Loader) buildFactorySet(tsSet []timeseries.TimeSerie, r *rand.Rand) (timeseries.FactorySet, error) {
	fs, err := timeseries.NewFactorySet(tsSet, r, l.StartDate, l.MockInterval)
	if err != nil {
		return timeseries.FactorySet{}, err
	}
	return fs, nil
}

func (l *Loader) getTimeSeries(wrSet []prompb.WriteRequest) ([]timeseries.TimeSerie, error) {
	counter := make(timeseries.Counter)
	for _, wr := range wrSet {
		err := timeseries.CountAllTimeSeries(&wr, &counter)
		if err != nil {
			return nil, err
		}
	}

	tsSet, err := counter.GetAllTimeSeries()
	log.Printf("Found %d time series", len(tsSet))
	if err != nil {
		return nil, err
	}

	return tsSet, nil
}

func (l *Loader) printStats(requests *[]prompb.WriteRequest) error {
	tableNameCounter := make(map[string]bool)

	for _, request := range *requests {
		err := timeseries.CountTableName(&request, &tableNameCounter)
		if err != nil {
			log.Fatalf("failed to count time series: %v", err)
		}
	}

	log.Printf("Found %d table names", len(tableNameCounter))

	return nil
}

func (l *Loader) getAllRemoteWriteRequest() ([]prompb.WriteRequest, error) {
	log.Printf("Parsing request file %s", l.TcpflowOutput)
	requests, err := parser.ParseHTTPRequests(l.TcpflowOutput)
	if err != nil {
		log.Printf("Error parsing request file %s: %v", l.TcpflowOutput, err)
		return nil, err
	}

	wrSet := []prompb.WriteRequest{}
	for _, request := range requests {
		wr, err := decode.Body(request.Body)
		if err != nil {
			log.Fatalf("failed to decode body: %v", err)
		}
		wrSet = append(wrSet, wr)
	}
	return wrSet, nil
}

func NewCommand() *cobra.Command {
	loader := &Loader{}

	var rootCmd = &cobra.Command{
		Use:   "loader",
		Short: "Loader is a tool to load data into a database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := loader.Run(cmd, args); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}

	rootCmd.Flags().StringP("url", "u", "http://localhost:4000/v1/prometheus/write?db=public", "The URL of the database")
	rootCmd.Flags().StringP("tcpflow-output", "t", "", "The path to the tcpflow output")
	rootCmd.Flags().IntP("timeseries-per-request", "r", 2000, "The number of timeseries per request")
	rootCmd.Flags().StringP("interval", "i", "1s", "The interval of the loading data")
	rootCmd.Flags().StringP("mock-interval", "m", "10s", "The interval of the mock data")
	rootCmd.Flags().IntP("seed", "s", 123456, "The seed for the random number generator")
	rootCmd.Flags().StringP("start-date", "", "2025-01-01T00:00:00Z", "The start date of the data")
	rootCmd.Flags().BoolP("dry-run", "d", false, "Dry run the loader")
	rootCmd.Flags().IntP("sample-scale", "", 1, "The scale of the samples")
	rootCmd.Flags().IntP("metric-scale", "", 1, "The scale of the metrics")
	rootCmd.Flags().IntP("label-scale", "", 1, "The scale of the labels")

	return rootCmd
}
