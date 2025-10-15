package samples

import (
	"math"

	"golang.org/x/exp/rand"
)

// LabelGenerator is an interface that represents a label generator.
type LabelGenerator interface {
	Next() string

	// All returns all the possible values
	All() []string

	NumCandidates() int
}

// FloatGenerator is an interface that represents a float generator.
type FloatGenerator interface {
	Next() float64
}

// MonoInc is a struct that represents a monotonically increasing float.
type MonoInc struct {
	step    int
	current int
}

// Next returns a monotonically increasing float.
func (m *MonoInc) Next() float64 {
	value := float64(m.current)
	m.current += m.step
	return value
}

// NewMonoInc creates a new MonoInc
func NewMonoInc(step int) *MonoInc {
	return &MonoInc{
		step:    step,
		current: 0,
	}
}

// RandomFloat is a struct that represents a random float.
type RandomFloat struct {
	lowerBound float64
	upperBound float64
}

// Next returns a random value from the float distribution.
func (r *RandomFloat) Next() float64 {
	return rand.Float64()*(r.upperBound-r.lowerBound) + r.lowerBound
}

// NewRandom creates a new RandomFloat
func NewRandom(lowerBound float64, upperBound float64) *RandomFloat {
	return &RandomFloat{
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
}

// RandomInt is a struct that represents a random integer.
type RandomInt struct {
	lowerBound int
	upperBound int
}

// Next returns a random integer.
func (r *RandomInt) Next() float64 {
	return float64(rand.Intn(r.upperBound-r.lowerBound) + r.lowerBound)
}

// NewRandomInt creates a new RandomInt
func NewRandomInt(lowerBound int, upperBound int) *RandomInt {
	return &RandomInt{lowerBound: lowerBound, upperBound: upperBound}
}

// ConstantString is a struct that represents a constant string.
type ConstantString struct {
	value string
}

func (c *ConstantString) NumCandidates() int {
	return 1
}

// Next returns the constant value.
func (c *ConstantString) Next() string {
	return c.value
}

// NewConstantString creates a new ConstantString
func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

// All returns the constant value.
func (c *ConstantString) All() []string {
	return []string{c.value}
}

// ConstantFloat is a struct that represents a constant float.
type ConstantFloat struct {
	value float64
}

// Next returns the constant value.
func (c *ConstantFloat) Next() float64 {
	return c.value
}

// NewConstantFloat creates a new ConstantFloat
func NewConstantFloat(value float64) *ConstantFloat {
	return &ConstantFloat{
		value: value,
	}
}

// Normal is a struct that represents a normal distribution.
type Normal struct {
	mean   float64
	stddev float64
}

// Next returns a random value from the normal distribution.
func (n *Normal) Next() float64 {
	return rand.NormFloat64()*n.stddev + n.mean
}

// NewNormal creates a new Normal
func NewNormal(mean float64, stddev float64) *Normal {
	return &Normal{mean: mean, stddev: stddev}
}

// Uniform is a struct that represents a uniform distribution.
type Uniform struct {
	lowerBound float64
	upperBound float64
}

// Next returns a random value from the uniform distribution.
func (u *Uniform) Next() float64 {
	return rand.Float64()*(u.upperBound-u.lowerBound) + u.lowerBound
}

// NewUniform creates a new Uniform
func NewUniform(lowerBound float64, upperBound float64) *Uniform {
	return &Uniform{lowerBound: lowerBound, upperBound: upperBound}
}

// Noisy is a struct that represents a noisy distribution.
type Noisy struct {
	current        float64
	maxFluctuation int
}

// Next returns a random value from the noisy distribution.
func (n *Noisy) Next() float64 {
	value := n.current
	n.current += rand.Float64()*float64(2*n.maxFluctuation) - float64(n.maxFluctuation)
	return value
}

// NewNoisy creates a new Noisy
func NewNoisy(maxFluctuation int) *Noisy {
	return &Noisy{maxFluctuation: maxFluctuation}
}

// Periodic is a struct that represents a periodic distribution.
type Periodic struct {
	period    int
	amplitude int
	bias      int
	current   float64
}

// Next returns a random value from the periodic distribution.
func (p *Periodic) Next() float64 {
	value := p.current
	p.current += float64(p.amplitude)*math.Sin(value/float64(p.period)) + float64(p.bias)
	return value
}

// NewPeriodic creates a new Periodic
func NewPeriodic(period int, amplitude int, bias int) *Periodic {
	return &Periodic{period: period, amplitude: amplitude, bias: bias}
}

// WeightedPreset is a struct that represents a weighted preset.
type WeightedPreset struct {
	preset      []PresetItem
	totalWeight int
}

func (w *WeightedPreset) NumCandidates() int {
	return len(w.preset)
}

// Next returns a random value from the preset based on the weighted distribution.
func (w *WeightedPreset) Next() string {
	for {
		r := rand.Intn(w.totalWeight)
		for _, item := range w.preset {
			if r < item.Weight {
				return item.Value
			}
			r -= item.Weight
		}
	}
}

// All returns all the possible values from the preset.
func (w *WeightedPreset) All() []string {
	all := make([]string, len(w.preset))
	for i, item := range w.preset {
		all[i] = item.Value
	}
	return all
}

// NewWeightedPreset creates a new WeightedPreset
func NewWeightedPreset(preset []PresetItem) *WeightedPreset {
	totalWeight := 0
	for _, item := range preset {
		totalWeight += item.Weight
	}
	return &WeightedPreset{preset: preset, totalWeight: totalWeight}
}
