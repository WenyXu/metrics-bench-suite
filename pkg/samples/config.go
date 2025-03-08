package samples

import (
	"fmt"
	"log"
	"strings"
)

// Config is a struct that represents a config for a metric.
type Config struct {
	Start    string  `yaml:"start"`
	End      string  `yaml:"end"`
	Interval int     `yaml:"interval"`
	Tags     []Tag   `yaml:"tags"`
	Fields   []Field `yaml:"fields"`
}

// Tag is a struct that represents a tag in a metric.
type Tag struct {
	Name string       `yaml:"name"`
	Type string       `yaml:"type"`
	Dist Distribution `yaml:"dist"`
}

// Field is a struct that represents a field in a metric.
type Field struct {
	Name string       `yaml:"name"`
	Type string       `yaml:"type"`
	Dist Distribution `yaml:"dist"`
}

// Distribution is a struct that represents a distribution for a tag or field.
type Distribution struct {
	Type           string       `yaml:"type"`
	LowerBound     *float64     `yaml:"lower_bound,omitempty"`
	UpperBound     *float64     `yaml:"upper_bound,omitempty"`
	Mean           *float64     `yaml:"mean,omitempty"`
	StdDev         *float64     `yaml:"stddev,omitempty"`
	Step           *int         `yaml:"step,omitempty"`
	Period         *int         `yaml:"period,omitempty"`
	Amplitude      *int         `yaml:"amplitude,omitempty"`
	Bias           *int         `yaml:"bias,omitempty"`
	Value          interface{}  `yaml:"value,omitempty"`
	MaxFluctuation *int         `yaml:"max_fluctuation,omitempty"`
	Preset         []PresetItem `yaml:"preset,omitempty"`
	Replica        *int         `yaml:"replica,omitempty"`
	ReplicaPrefix  *string      `yaml:"replica_prefix,omitempty"`
}

// PresetItem is a struct that represents an item in a preset.
type PresetItem struct {
	Value  string `yaml:"value"`
	Weight int    `yaml:"weight"`
}

// LabelDistributionType is a type that represents a distribution type for a label.
type LabelDistributionType int

const (
	weightedPresetType LabelDistributionType = iota
	constantStringType
	replicaType
)

var distributionTypeMap = map[string]LabelDistributionType{
	"weighted_preset": weightedPresetType,
	"constant_string": constantStringType,
	"replica_string":  replicaType,
}

// parseDistributionTypeFromString parses a distribution type from a string
func parseDistributionTypeFromString(s string) (LabelDistributionType, error) {
	s = strings.ToLower(s)
	if v, exist := distributionTypeMap[s]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("invalid distribution type: %s", s)
}

// Len returns the length of the distribution
func (d Distribution) Len() int {
	distributionType, err := parseDistributionTypeFromString(d.Type)
	if err != nil {
		log.Fatalf("invalid distribution type: %s", d.Type)
	}

	switch distributionType {
	case weightedPresetType:
		return len(d.Preset)
	case constantStringType:
		return 1
	case replicaType:
		return *d.Replica
	}

	return 0
}

// LabelGenerator returns a label generator for the distribution
func (d Distribution) LabelGenerator() LabelGenerator {
	distributionType, err := parseDistributionTypeFromString(d.Type)
	if err != nil {
		log.Fatalf("invalid distribution type: %s", d.Type)
	}

	switch distributionType {
	case weightedPresetType:
		return NewWeightedPreset(d.Preset)
	case constantStringType:
		return NewConstantString(d.Value.(string))
	case replicaType:
		preset := make([]PresetItem, *d.Replica)
		for i := 0; i < *d.Replica; i++ {
			preset[i] = PresetItem{
				Value:  fmt.Sprintf("%s%d", *d.ReplicaPrefix, i),
				Weight: 1,
			}
		}
		return NewWeightedPreset(preset)
	}

	return nil
}

// FieldDistributionType is a type that represents a distribution type for a field.
type FieldDistributionType int

const (
	monoIncType FieldDistributionType = iota
	normalType
	randomFloatType
	randomIntType
	constantFloatType
	uniformType
	noisyType
	periodicType
)

var fieldDistributionTypeMap = map[string]FieldDistributionType{
	"mono_inc":       monoIncType,
	"normal":         normalType,
	"random_float":   randomFloatType,
	"random_int":     randomIntType,
	"constant_float": constantFloatType,
	"uniform":        uniformType,
	"noisy":          noisyType,
	"periodic":       periodicType,
}

func parseFieldDistributionTypeFromString(s string) (FieldDistributionType, error) {
	s = strings.ToLower(s)
	if v, exist := fieldDistributionTypeMap[s]; exist {
		return v, nil
	}
	return 0, fmt.Errorf("invalid field distribution type: %s", s)
}

// FieldGenerator returns a float generator for the distribution
func (d Distribution) FieldGenerator() FloatGenerator {
	distributionType, err := parseFieldDistributionTypeFromString(d.Type)
	if err != nil {
		log.Fatalf("invalid field distribution type: %s", d.Type)
	}

	switch distributionType {
	case monoIncType:
		return NewMonoInc(*d.Step)
	case normalType:
		return NewNormal(*d.Mean, *d.StdDev)
	case randomIntType:
		return NewRandomInt(int(*d.LowerBound), int(*d.UpperBound))
	case randomFloatType:
		return NewRandom(*d.LowerBound, *d.UpperBound)
	case constantFloatType:
		return NewConstantFloat(d.Value.(float64))
	case uniformType:
		return NewUniform(*d.LowerBound, *d.UpperBound)
	case noisyType:
		return NewNoisy(*d.MaxFluctuation)
	case periodicType:
		return NewPeriodic(*d.Period, *d.Amplitude, *d.Bias)
	}

	return nil
}
