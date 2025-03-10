package main

import (
	"encoding/json"
	"fmt"
	"log"

	"time"

	"metrics-bench-suite/pkg/http"
	"metrics-bench-suite/pkg/samples"

	"github.com/prometheus/prometheus/prompb"
	"github.com/spf13/cobra"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

// SampleGenerator is a struct that generates samples from a config file
type SampleGenerator struct {
	ConfigPath     string
	Interval       time.Duration
	StartDate      time.Time
	EndDate        time.Time
	Seed           int
	Output         string
	RemoteWriteURL string
	Database       string
	Table          string
}

type metric struct {
	Name   string
	Series []map[string]string
	Fields []samples.FloatGenerator
}

func (s *SampleGenerator) run(cmd *cobra.Command, args []string) error {
	var err error
	intervalStr, _ := cmd.Flags().GetString("interval")
	initialDateStr, _ := cmd.Flags().GetString("start-date")
	endDateStr, _ := cmd.Flags().GetString("end-date")
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
	s.Seed, err = cmd.Flags().GetInt("seed")
	if err != nil {
		return err
	}
	s.Output, err = cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	s.RemoteWriteURL, err = cmd.Flags().GetString("remote-write-url")
	if err != nil {
		return err
	}
	s.Database, err = cmd.Flags().GetString("database")
	if err != nil {
		return err
	}
	s.Table, err = cmd.Flags().GetString("table")
	if err != nil {
		return err
	}
	log.Printf("Start date: %s", s.StartDate)
	log.Printf("End date: %s", s.EndDate)
	log.Printf("Interval: %s", s.Interval)
	log.Printf("Seed: %d", s.Seed)
	log.Printf("Config path: %s", s.ConfigPath)
	log.Printf("Output: %s", s.Output)

	fileConfigs, err := samples.WalkAndParseConfig(s.ConfigPath)
	if err != nil {
		return err
	}
	if len(fileConfigs) == 0 {
		return fmt.Errorf("no config files found")
	}

	log.Printf("Generating metrics...")
	totalCount := 0
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
		tagSet := samples.TagSetPermutation(labels, &totalCount)
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

	log.Printf("total time series: %d, for %d metrics", totalCount, len(metrics))

	if s.RemoteWriteURL != "" {
		log.Printf("Sending metrics to remote write...")
		wr := convertToRemoteWriteRequest(metrics, s.StartDate, s.EndDate, s.Interval)
		err = http.NewRequester(s.RemoteWriteURL).Send(wr)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Saving metrics to file...")
		parquetFile, err := local.NewLocalFileWriter(s.fileName())
		if err != nil {
			return err
		}
		parquetWriter, err := newParquetWriter(parquetFile)
		if err != nil {
			return err
		}

		for _, metric := range metrics {
			tss := convertMetricToTimeSeries(metric, s.StartDate, s.EndDate, s.Interval)
			for _, ts := range tss {
				json, err := json.Marshal(ts)
				if err != nil {
					return err
				}
				err = parquetWriter.Write(row{Value: string(json)})
				if err != nil {
					return err
				}
			}
		}

		if err := parquetWriter.WriteStop(); err != nil {
			return err
		}
		if err := parquetFile.Close(); err != nil {
			return err
		}

	}

	return nil
}

func newParquetWriter(parquetFile source.ParquetFile) (*writer.ParquetWriter, error) {
	parquetWriter, err := writer.NewParquetWriter(parquetFile, new(row), 1)
	if err != nil {
		return nil, err
	}

	parquetWriter.RowGroupSize = 102400
	parquetWriter.PageSize = 1024 * 1024
	parquetWriter.CompressionType = parquet.CompressionCodec_ZSTD
	return parquetWriter, nil
}

type row struct {
	Value string `parquet:"name=value, type=BYTE_ARRAY, convertedtype=UTF8"`
}

func convertMetricToTimeSeries(metric metric, start time.Time, end time.Time, interval time.Duration) []prompb.TimeSeries {
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
			ts.Labels = append(ts.Labels, prompb.Label{
				Name:  k,
				Value: v,
			})
		}

		generator := metric.Fields[i]
		current := start
		for current.Before(end) {
			ts.Samples = append(ts.Samples, prompb.Sample{
				Value:     generator.Next(),
				Timestamp: current.UnixMilli(),
			})
			current = current.Add(interval)
		}
		tsSet = append(tsSet, ts)
	}

	return tsSet
}

func convertToRemoteWriteRequest(metrics []metric, start time.Time, end time.Time, interval time.Duration) prompb.WriteRequest {
	tsSet := make([]prompb.TimeSeries, 0)
	for _, metric := range metrics {
		tsSet = append(tsSet, convertMetricToTimeSeries(metric, start, end, interval)...)
	}

	return prompb.WriteRequest{
		Timeseries: tsSet,
	}
}

var timeFormat = "2006-01-02T150405Z"

func (s *SampleGenerator) fileName() string {
	startDate := s.StartDate.Format(timeFormat)
	endDate := s.EndDate.Format(timeFormat)
	return fmt.Sprintf("%s.%s-[%s-%s].%s", s.Database, s.Table, startDate, endDate, "parquet")
}

func main() {
	sampleGenerator := &SampleGenerator{}

	var rootCmd = &cobra.Command{
		Use:   "sample_generator",
		Short: "SampleGenerator is a tool to generate timeseries samples from a config file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sampleGenerator.run(cmd, args); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}

	rootCmd.Flags().StringP("database", "d", "public", "The database name")
	rootCmd.Flags().StringP("table", "t", "greptime_physical_table", "The table name")
	rootCmd.Flags().StringP("config", "c", "", "The path to the config file")
	rootCmd.Flags().StringP("interval", "i", "30s", "The interval of the loading data")
	rootCmd.Flags().StringP("start-date", "", "2025-01-01T00:00:00Z", "The start date of the data")
	rootCmd.Flags().StringP("end-date", "", "2025-01-01T00:01:00Z", "The end date of the data")
	rootCmd.Flags().IntP("seed", "s", 123456, "The seed for the random number generator")
	rootCmd.Flags().StringP("output", "o", "output.json", "The output file")
	rootCmd.Flags().StringP("remote-write-url", "u", "", "The remote write url")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
