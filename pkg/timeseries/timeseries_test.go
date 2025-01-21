package timeseries

import (
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

func TestConvertTimeSeriesToLabels(t *testing.T) {
	wr := &prompb.WriteRequest{
		Timeseries: []prompb.TimeSeries{
			{
				Labels: []prompb.Label{
					{Name: "__name__", Value: "test_table"},
					{Name: "pod_name", Value: "test_pod"},
					{Name: "namespace", Value: "test_namespace"},
					{Name: "container_name", Value: "test_container"},
				},
			}, {
				Labels: []prompb.Label{
					{Name: "__name__", Value: "table_1"},
					{Name: "pod_name", Value: "test_pod_1"},
					{Name: "namespace", Value: "test_namespace_1"},
					{Name: "container_name", Value: "test_container_1"},
				},
			},
		},
	}
	counter := make(Counter)
	err := CountAllTimeSeries(wr, &counter)
	assert.Nil(t, err)

	ts, err := counter.GetAllTimeSeries()
	if err != nil {
		t.Fatalf("failed to get all time series: %v", err)
	}

	assert.Equal(t, len(ts), 2)

	_, err = counter.GetAllTimeSeries()
	assert.Nil(t, err)

}
