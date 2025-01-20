package timeseries

import (
	"fmt"
	"metrics-bench-suite/pkg/utils/decode"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
)

func TestFindTimeSeries(t *testing.T) {
	input := &prompb.TimeSeries{
		Labels: []prompb.Label{
			{Name: "table_name", Value: "test_table"},
			{Name: "pod_name", Value: "test_pod"},
			{Name: "namespace", Value: "test_namespace"},
			{Name: "container_name", Value: "test_container"},
		},
	}

	tableName, ts, err := ExtractTimeSeries(input)
	if err != nil {
		t.Fatalf("failed to find time series: %v", err)
	}

	assert.Equal(t, tableName, "test_table")
	assert.Equal(t, ts.String(), "pod_name=test_pod, namespace=test_namespace, container_name=test_container")
}

func TestParseTimeSeries(t *testing.T) {
	ts := "pod_name=test_pod, namespace=test_namespace, container_name=test_container"
	timeseries, err := ParseTimeSeries(ts)
	if err != nil {
		t.Fatalf("failed to parse time series: %v", err)
	}
	expected := []Label{
		{Name: "pod_name", Value: "test_pod"},
		{Name: "namespace", Value: "test_namespace"},
		{Name: "container_name", Value: "test_container"},
	}

	assert.Equal(t, expected, timeseries.Labels)
}

func TestCountTimeSeries(t *testing.T) {
	root := "../../assets/"
	assets, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	tableNameCounter := make(map[string]bool)
	timeSeriesCounter := make(map[string]bool)

	for _, entry := range assets {
		// Check if the entry is a file
		if entry.IsDir() {
			fmt.Printf("Directory: %s\n", entry.Name())
		} else {
			data, err := os.ReadFile(filepath.Join(root, entry.Name()))
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			compressed, err := decode.Base64(string(data))
			assert.Nil(t, err)
			wr, err := decode.Body(compressed)
			err = CountTimeSeries(&wr, &tableNameCounter, &timeSeriesCounter)
			if err != nil {
				t.Fatalf("failed to count time series: %v", err)
			}
		}
	}

	println("tableNameCounter: ", len(tableNameCounter))
	println("timeSeriesCounter: ", len(timeSeriesCounter))
}
