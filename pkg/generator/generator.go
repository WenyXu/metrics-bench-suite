package generator

// Value represents a value with a date and a value
type Value struct {
	Timestamp int64
	Value     float64
}

// Generator is an interface that represents a generator
type Generator interface {
	Generate() (Value, error)
}
