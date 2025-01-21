package generator

import (
	"time"
)

// LinearTrend struct represents a linear trend factor, which can either be simple or feature-based
type LinearTrend struct {
	Coefficient float64 // The coefficient (coef) of the linear trend
	Offset      float64 // The offset (b) of the linear trend
	Step        time.Duration
	InitialDate time.Time
	CurrentDate time.Time
}

// NewLinearTrend creates a new instance of LinearTrend and validates inputs
func NewLinearTrend(coef, offset float64, initialDate time.Time, step time.Duration) (*LinearTrend, error) {

	// Create and return a new LinearTrend instance
	return &LinearTrend{
		Coefficient: coef,
		Offset:      offset,
		Step:        step,
		InitialDate: initialDate,
		CurrentDate: initialDate,
	}, nil
}

// Generate generates the linear trend data over a given date range
func (lt *LinearTrend) Generate() (Value, error) {
	lt.CurrentDate = lt.CurrentDate.Add(lt.Step)
	// Calculate the number of days since the start date
	days := float64(lt.CurrentDate.Sub(lt.InitialDate).Hours() / 24)
	// Calculate y = ax + b (with coef as a, offset as b) using the previous value
	value := lt.Coefficient*days + 1 + lt.Offset
	return Value{Timestamp: lt.CurrentDate.Unix(), Value: value}, nil
}
