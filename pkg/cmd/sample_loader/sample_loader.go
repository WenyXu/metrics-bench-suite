package sample_loader

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"metrics-bench-suite/pkg/http"
	"metrics-bench-suite/pkg/samples"
	"sync"
	"time"

	"github.com/prometheus/prometheus/prompb"
	"github.com/spf13/cobra"
)

// SampleLoader is a tool that generate samples from config files and send them to the remote write endpoint.
type SampleLoader struct {
	ConfigPath     string
	RemoteWriteURL string
	StartDate      time.Time
	EndDate        time.Time
	Interval       time.Duration
	Seed           int
	MaxSamples     int
	TickInterval   time.Duration
	Workers        int
	Infinite       bool
	TagsPickRate   float32
	TablePickCount uint64
}

type metric struct {
	Name   string
	Series []map[string]string
	Fields []samples.FloatGenerator
}

func (s *SampleLoader) run(cmd *cobra.Command, _ []string) error {
	var err error
	intervalStr, _ := cmd.Flags().GetString("interval")
	initialDateStr, _ := cmd.Flags().GetString("start-date")
	endDateStr, _ := cmd.Flags().GetString("end-date")
	tickIntervalStr, _ := cmd.Flags().GetString("tick-interval")
	s.Interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}
	s.ConfigPath, err = cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	s.StartDate, err = time.Parse(time.RFC3339, initialDateStr)
	if err != nil {
		return err
	}
	s.EndDate, err = time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return err
	}
	s.RemoteWriteURL, err = cmd.Flags().GetString("remote-write-url")
	if err != nil {
		return err
	}
	s.MaxSamples, err = cmd.Flags().GetInt("max-samples")
	if err != nil {
		return err
	}
	s.TickInterval, err = time.ParseDuration(tickIntervalStr)
	if err != nil {
		return err
	}
	s.Workers, err = cmd.Flags().GetInt("workers")
	if err != nil {
		return err
	}
	s.Infinite, err = cmd.Flags().GetBool("infinite")
	if err != nil {
		return err
	}
	s.TagsPickRate, err = cmd.Flags().GetFloat32("tags-pick-rate")
	if err != nil {
		return err
	}
	s.TablePickCount, err = cmd.Flags().GetUint64("table-pick-count")
	if err != nil {
		return err
	}
	log.Printf("Start date: %s", s.StartDate)
	log.Printf("End date: %s", s.EndDate)
	log.Printf("Interval: %s", s.Interval)
	log.Printf("Tick interval: %s", s.TickInterval)
	log.Printf("Config path: %s", s.ConfigPath)
	log.Printf("Tags pick rate: %f", s.TagsPickRate)
	log.Printf("Table pick rate: %d", s.TablePickCount)

	fileConfigs, err := samples.WalkAndParseConfigWithMaxFileCount(s.ConfigPath, s.TablePickCount)
	if err != nil {
		return err
	}
	if len(fileConfigs) == 0 {
		return fmt.Errorf("no config files found")
	}

	log.Printf("Generating metrics...")
	totalSeries := 0
	metrics := make([]metric, len(fileConfigs))
	for i, fileConfig := range fileConfigs {
		labels := make([]samples.Label, 0)
		for _, tag := range fileConfig.Config.Tags {
			values := tag.Dist.LabelGenerator().All()
			labels = append(labels, samples.Label{
				Name:   tag.Name,
				Values: values,
			})
		}

		log.Printf("Process %s", fileConfig.Name)
		tagSet := samples.TagSetPermutation(labels, &totalSeries)
		fields := make([]samples.FloatGenerator, 0)
		field := fileConfig.Config.Fields[0]
		for range len(tagSet) {
			fields = append(fields, field.Dist.FieldGenerator())
		}
		metrics[i] = metric{
			Name:   fileConfig.Name,
			Series: tagSet,
			Fields: fields,
		}
	}

	log.Printf("total time series: %d, for %d metrics", totalSeries, len(metrics))

	ticker := time.NewTicker(s.TickInterval)
	defer ticker.Stop()

	requestChan := make(chan prompb.WriteRequest, s.Workers)

	wg := sync.WaitGroup{}
	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go worker(i, s.RemoteWriteURL, requestChan, &wg)
	}

	current := s.StartDate
	if s.Infinite {
		current = time.Now()
	}

	for range ticker.C {
		log.Printf("Generating samples for %s", current)
		convertToRemoteWriteRequests(metrics, current, s.MaxSamples, requestChan, s.TagsPickRate)
		current = current.Add(s.Interval)
		if !s.Infinite {
			if current.After(s.EndDate) {
				log.Printf("End date reached, stopping")
				break
			}
		}
	}

	close(requestChan)
	wg.Wait()

	return nil
}

func worker(id int, url string, request <-chan prompb.WriteRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range request {
		now := time.Now()
		r := http.NewRequester(url)
		err := r.Send(request)
		if err != nil {
			log.Printf("worker %d failed to send write request: %v", id, err)
		}
		log.Printf("worker %d sent request in %s", id, time.Since(now))
	}
}

func convertMetricToTimeSeries(metric metric, current time.Time, pickRate float32) []prompb.TimeSeries {
	tsSet := make([]prompb.TimeSeries, 0)
	for i, label := range metric.Series {
		ts := prompb.TimeSeries{
			Labels:  make([]prompb.Label, 0),
			Samples: make([]prompb.Sample, 0),
		}
		ts.Labels = append(ts.Labels, prompb.Label{
			Name:  "__name__",
			Value: metric.Name,
		})
		for k, v := range label {
			if pickRate < 1.0 {
				if rand.Float32() > pickRate {
					continue
				}
			}
			ts.Labels = append(ts.Labels, prompb.Label{
				Name:  k,
				Value: v,
			})
		}

		generator := metric.Fields[i]
		ts.Samples = append(ts.Samples, prompb.Sample{
			Value:     generator.Next(),
			Timestamp: current.UnixMilli(),
		})
		tsSet = append(tsSet, ts)
	}

	return tsSet
}

func convertToRemoteWriteRequests(metrics []metric, current time.Time, maxSamples int, requestChan chan<- prompb.WriteRequest, pickRate float32) {

	tsSet := make([]prompb.TimeSeries, 0)
	for _, metric := range metrics {
		tsSet = append(tsSet, convertMetricToTimeSeries(metric, current, pickRate)...)
	}
	for len(tsSet) > 0 {
		if len(tsSet) > maxSamples {
			requestChan <- prompb.WriteRequest{
				Timeseries: tsSet[:maxSamples],
			}
			tsSet = tsSet[maxSamples:]
		} else {
			requestChan <- prompb.WriteRequest{
				Timeseries: tsSet,
			}
			tsSet = make([]prompb.TimeSeries, 0)
		}
	}
}

func NewCommand() *cobra.Command {
	sampleLoader := &SampleLoader{}

	var rootCmd = &cobra.Command{
		Use:   "sample_loader",
		Short: "SampleLoader is a tool to load samples from a file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sampleLoader.run(cmd, args); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}

	rootCmd.Flags().StringP("config", "c", "", "The path to the config file")
	rootCmd.Flags().StringP("remote-write-url", "u", "", "The remote write url")
	rootCmd.Flags().StringP("start-date", "", "2025-01-01T00:00:00Z", "The start date of the data")
	rootCmd.Flags().StringP("end-date", "", "2025-01-01T00:01:00Z", "The end date of the data")
	rootCmd.Flags().StringP("interval", "", "30s", "The interval of the data")
	rootCmd.Flags().IntP("max-samples", "s", 20000, "The max number of metrics to load")
	rootCmd.Flags().StringP("tick-interval", "t", "30s", "The interval of the requests")
	rootCmd.Flags().IntP("workers", "w", 1, "The number of workers to send requests")
	rootCmd.Flags().BoolP("infinite", "i", false, "Run indefinitely")
	rootCmd.Flags().Float32P("tags-pick-rate", "p", 1.0, "The rate of the pick tags")
	rootCmd.Flags().Uint64P("table-pick-count", "n", math.MaxUint64, "The number of tables to pick from")

	return rootCmd
}
