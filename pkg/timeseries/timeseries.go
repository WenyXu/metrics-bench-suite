package timeseries

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/prompb"
)

// TimeSerie represents a time series with labels
type TimeSerie struct {
	Labels []Label
}

// Label represents a label with a name and value
type Label struct {
	Name  string
	Value string
}

// String implements the Stringer interface for TimeSeries.
func (ts TimeSerie) String() string {
	var parts []string
	for _, label := range ts.Labels {
		parts = append(parts, fmt.Sprintf("%s=%s", label.Name, label.Value))
	}
	return strings.Join(parts, ", ")
}

// ParseTimeSeries parses a serialized string into a TimeSeries object.
func ParseTimeSeries(input string) (TimeSerie, error) {
	// Initialize an empty TimeSeries
	ts := TimeSerie{}

	// Split the string by ", " to separate labels
	labelPairs := strings.Split(input, ", ")
	for _, pair := range labelPairs {
		// Split each pair by ": " to get the key and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return ts, errors.New("invalid format: each label must be in 'key=value' format")
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
func ExtractTimeSeries(ts *prompb.TimeSeries) (string, TimeSerie, error) {
	label, err := ExtractFirstLabel(ts)
	if err != nil {
		return "", TimeSerie{}, err
	}

	tableName := label.Value

	var labels []Label
	for _, label := range ts.Labels[1:] {
		labels = append(labels, Label{
			Name:  label.Name,
			Value: label.Value,
		})
	}

	return tableName, TimeSerie{Labels: labels}, nil
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

// Counter is a map of TimeSeries to a boolean
type Counter map[string]TimeSerie

// GetAllTimeSeries returns all the time series in the counter
func (c *Counter) GetAllTimeSeries() ([]TimeSerie, error) {
	output := []TimeSerie{}
	for _, value := range *c {
		output = append(output, value)
	}
	return output, nil
}

// ConvertTimeSeriesToLabels converts a TimeSeries to a list of labels
func ConvertTimeSeriesToLabels(ts TimeSerie) []prompb.Label {
	labels := []prompb.Label{}
	for _, label := range ts.Labels {
		labels = append(labels, prompb.Label{Name: label.Name, Value: label.Value})
	}
	return labels
}

// ConvertLabelsToTimeSerie converts a list of labels to a TimeSeries
func ConvertLabelsToTimeSerie(labels []prompb.Label) TimeSerie {
	ts := TimeSerie{}
	for _, label := range labels {
		ts.Labels = append(ts.Labels, Label{Name: label.Name, Value: label.Value})
	}
	return ts
}

// CountAllTimeSeries counts the number of time series per table
func CountAllTimeSeries(wr *prompb.WriteRequest, tableNameCounter *Counter) error {
	for _, ts := range wr.Timeseries {
		timeserie := ConvertLabelsToTimeSerie(ts.Labels)
		if _, ok := (*tableNameCounter)[timeserie.String()]; !ok {
			(*tableNameCounter)[timeserie.String()] = timeserie
		}
	}
	return nil
}
