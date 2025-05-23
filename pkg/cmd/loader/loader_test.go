package loader

import (
	"testing"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
)

func TestScaleMetrics(t *testing.T) {
	tsSet := []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "total"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
	}

	expected := []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "total"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric_1"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "total_1"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
	}
	scaled := ScaleMetrics(tsSet, 2)
	assert.Equal(t, len(scaled), 2*len(tsSet))
	assert.Equal(t, expected, scaled)

}

func TestScaleLabels(t *testing.T) {
	tsSet := []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "job", Value: "job"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
	}

	expected := []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "host", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "job", Value: "job"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "hosta", Value: "host"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "metric"},
				{Name: "joba", Value: "job"},
			},
			Samples: []prompb.Sample{
				{Value: 1, Timestamp: 1},
			},
		},
	}

	scaled := ScaleLabels(tsSet, 2)
	assert.Equal(t, expected, scaled)
}
