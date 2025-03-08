package metric

import (
	"math"

	"golang.org/x/exp/rand"
)

type MonoInc struct {
	step    int
	current int
}

func (m *MonoInc) Next() float64 {
	value := float64(m.current)
	m.current += m.step
	return value
}

func NewMonoInc(step int) *MonoInc {
	return &MonoInc{
		step:    step,
		current: 0,
	}
}

type RandomFloat struct {
	lowerBound float64
	upperBound float64
}

func (r *RandomFloat) Next() float64 {
	return rand.Float64()*(r.upperBound-r.lowerBound) + r.lowerBound
}

func NewRandom(lowerBound float64, upperBound float64) *RandomFloat {
	return &RandomFloat{
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
}

type RandomInt struct {
	lowerBound int
	upperBound int
}

func (r *RandomInt) Next() int {
	return rand.Intn(r.upperBound-r.lowerBound) + r.lowerBound
}

func NewRandomInt(lowerBound int, upperBound int) *RandomInt {
	return &RandomInt{lowerBound: lowerBound, upperBound: upperBound}
}

type ConstantString struct {
	value string
}

func (c *ConstantString) Next() string {
	return c.value
}

func NewConstantString(value string) *ConstantString {
	return &ConstantString{
		value: value,
	}
}

type ConstantFloat struct {
	value float64
}

func (c *ConstantFloat) Next() float64 {
	return c.value
}

func NewConstantFloat(value float64) *ConstantFloat {
	return &ConstantFloat{
		value: value,
	}
}

type Normal struct {
	mean   float64
	stddev float64
}

func (n *Normal) Next() float64 {
	return rand.NormFloat64()*n.stddev + n.mean
}

func NewNormal(mean float64, stddev float64) *Normal {
	return &Normal{mean: mean, stddev: stddev}
}

type Uniform struct {
	lowerBound float64
	upperBound float64
}

func (u *Uniform) Next() float64 {
	return rand.Float64()*(u.upperBound-u.lowerBound) + u.lowerBound
}

func NewUniform(lowerBound float64, upperBound float64) *Uniform {
	return &Uniform{lowerBound: lowerBound, upperBound: upperBound}
}

type Noisy struct {
	current         float64
	max_fluctuation int
}

func (n *Noisy) Next() float64 {
	value := n.current
	n.current += rand.Float64()*float64(2*n.max_fluctuation) - float64(n.max_fluctuation)
	return value
}

func NewNoisy(max_fluctuation int) *Noisy {
	return &Noisy{max_fluctuation: max_fluctuation}
}

type Periodic struct {
	period    int
	amplitude int
	bias      int
	current   float64
}

func (p *Periodic) Next() float64 {
	value := p.current
	p.current += float64(p.amplitude)*math.Sin(value/float64(p.period)) + float64(p.bias)
	return value
}

func NewPeriodic(period int, amplitude int, bias int) *Periodic {
	return &Periodic{period: period, amplitude: amplitude, bias: bias}
}

type WeightedPreset struct {
	preset      []PresetItem
	totalWeight int
}

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

func NewWeightedPreset(preset []PresetItem) *WeightedPreset {
	totalWeight := 0
	for _, item := range preset {
		totalWeight += item.Weight
	}
	return &WeightedPreset{preset: preset, totalWeight: totalWeight}
}
