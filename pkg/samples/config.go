package metric

type Config struct {
	Start    string  `yaml:"start"`
	End      string  `yaml:"end"`
	Interval int     `yaml:"interval"`
	Tags     []Tag   `yaml:"tags"`
	Fields   []Field `yaml:"fields"`
}

type Tag struct {
	Name string       `yaml:"name"`
	Type string       `yaml:"type"`
	Dist Distribution `yaml:"dist"`
}

type Field struct {
	Name string       `yaml:"name"`
	Type string       `yaml:"type"`
	Dist Distribution `yaml:"dist"`
}

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
	Value          *string      `yaml:"value,omitempty"`
	MaxFluctuation *int         `yaml:"max_fluctuation,omitempty"`
	Preset         []PresetItem `yaml:"preset,omitempty"`
	Replica        *int         `yaml:"replica,omitempty"`
	ReplicaPrefix  *string      `yaml:"replica_prefix,omitempty"`
}

type PresetItem struct {
	Value  string `yaml:"value"`
	Weight int    `yaml:"weight"`
}
