package timeseries

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/prompb"
)

// TimeSeries represents a time series with labels
type TimeSeries struct {
	Labels []Label
}

// Label represents a label with a name and value
type Label struct {
	Name  string
	Value string
}

// String implements the Stringer interface for TimeSeries.
func (ts TimeSeries) String() string {
	var parts []string
	for _, label := range ts.Labels {
		parts = append(parts, fmt.Sprintf("%s=%s", label.Name, label.Value))
	}
	return strings.Join(parts, ", ")
}

// ParseTimeSeries parses a serialized string into a TimeSeries object.
func ParseTimeSeries(input string) (TimeSeries, error) {
	// Initialize an empty TimeSeries
	ts := TimeSeries{}

	// Split the string by ", " to separate labels
	labelPairs := strings.Split(input, ", ")
	for _, pair := range labelPairs {
		// Split each pair by ": " to get the key and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return ts, errors.New("invalid format: each label must be in 'key: value' format")
		}

		// Trim spaces around the key and value and add to the TimeSeries
		label := Label{
			Name:  strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		}
		ts.Labels = append(ts.Labels, label)
	}

	return ts, nil
}

// ExtractFirstLabel extracts the first label from the TimeSeries.
func ExtractFirstLabel(ts *prompb.TimeSeries) (prompb.Label, error) {
	if len(ts.Labels) == 0 {
		return prompb.Label{}, fmt.Errorf("the first timeseries has no labels")
	}
	return ts.Labels[0], nil
}

// ExtractTimeSeries extracts the table name and the time series from the given TimeSeries.
func ExtractTimeSeries(ts *prompb.TimeSeries) (string, TimeSeries, error) {
	label, err := ExtractFirstLabel(ts)
	if err != nil {
		return "", TimeSeries{}, err
	}

	tableName := label.Value

	var labels []Label
	for _, label := range ts.Labels[1:] {
		labels = append(labels, Label{
			Name:  label.Name,
			Value: label.Value,
		})
	}

	return tableName, TimeSeries{Labels: labels}, nil
}

// CountTimeSeries counts the number of time series in the WriteRequest.
func CountTimeSeries(wr *prompb.WriteRequest, tableNameCounter *map[string]bool, timeSeriesCounter *map[string]bool) error {
	for _, ts := range wr.Timeseries {
		tableName, timeseries, err := ExtractTimeSeries(&ts)
		if err != nil {
			return err
		}

		(*tableNameCounter)[tableName] = true
		(*timeSeriesCounter)[timeseries.String()] = true
	}

	return nil
}
