package timeseries

import (
	"math/rand"
	"metrics-bench-suite/pkg/generator"
	"time"

	"github.com/prometheus/prometheus/prompb"
)

// Factory represents a factory for a single TimeSeries
type Factory struct {
	labels    []prompb.Label
	generator generator.Generator
}

// Generate generates a list of labels and a list of values
func (tsf *Factory) Generate() ([]prompb.Label, generator.Value, error) {
	value, err := tsf.generator.Generate()
	if err != nil {
		return nil, generator.Value{}, err
	}
	return tsf.labels, value, nil
}

// FactorySet represents a set of TimeSeriesFactory
type FactorySet struct {
	timeseriesSet []Factory
}

// NewFactorySet creates a new FactorySet
func NewFactorySet(timeseries []TimeSerie, r *rand.Rand, initialDate time.Time, step time.Duration) (FactorySet, error) {
	timeseriesSet := []Factory{}
	for _, ts := range timeseries {
		cof := r.Int63n(100)
		offset := r.Int31n(10000)
		generator, err := generator.NewLinearTrend(float64(cof), float64(offset), initialDate, step)
		if err != nil {
			return FactorySet{}, err
		}
		timeseriesSet = append(timeseriesSet, Factory{labels: ConvertTimeSeriesToLabels(ts), generator: generator})
	}

	return FactorySet{timeseriesSet: timeseriesSet}, nil
}

// Sample represents a sample of a TimeSeries and its Value
type Sample struct {
	TimeSeries TimeSerie
	Value      generator.Value
}

// Generate generates a list of TimeSeries and a list of Values
func (fs *FactorySet) Generate(scale int) ([]prompb.TimeSeries, error) {
	tsSet := []prompb.TimeSeries{}

	for i := 0; i < scale; i++ {
		for _, tsf := range fs.timeseriesSet {
			labels, value, err := tsf.Generate()
			if err != nil {
				return nil, err
			}
			tsSet = append(tsSet, prompb.TimeSeries{Labels: labels, Samples: []prompb.Sample{{Value: value.Value, Timestamp: value.Timestamp}}})
		}
	}
	return tsSet, nil
}
