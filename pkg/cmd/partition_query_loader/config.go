package partition_query_loader

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	defaultMySQLHost        = "127.0.0.1"
	defaultMySQLPort        = 4002
	defaultDatabases        = "public"
	defaultConcurrency      = 1
	defaultDuration         = time.Duration(0)
	defaultPeriod           = 30 * time.Second
	defaultPeriodMin        = 15 * time.Second
	defaultPeriodMax        = 45 * time.Second
	defaultTimeRange        = 5 * time.Minute
	defaultTimeRangeMax     = 15 * time.Minute
	defaultHotPartitions    = 3
	defaultHotShare         = 0.85
	defaultStatsInterval    = 15 * time.Second
	defaultProfile          = "sustained"
	defaultOutputFormat     = "ndjson"
	defaultDiscoveryTimeout = 30 * time.Second
	defaultDialTimeout      = 10 * time.Second
	defaultReadTimeout      = 30 * time.Second
	defaultWriteTimeout     = 30 * time.Second
)

// Config holds the validated partition query loader settings.
type Config struct {
	MySQLHost         string
	MySQLPort         int
	MySQLUser         string
	MySQLPassword     string
	Databases         []string
	TablesPerDatabase int
	RandomSeed        *int64
	Concurrency       int
	Duration          time.Duration
	Period            time.Duration
	PeriodMin         time.Duration
	PeriodMax         time.Duration
	Jitter            time.Duration
	TimeRange         time.Duration
	TimeRangeMin      time.Duration
	TimeRangeMax      time.Duration
	HotPartitions     int
	HotShare          float64
	StatsInterval     time.Duration
	Profile           string
	OutputFormat      string
	DryRun            bool
	DiscoveryTimeout  time.Duration
	DialTimeout       time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

// DefaultConfig returns the loader configuration used when flags are omitted.
func DefaultConfig() Config {
	return Config{
		MySQLHost:         defaultMySQLHost,
		MySQLPort:         defaultMySQLPort,
		Databases:         []string{defaultDatabases},
		TablesPerDatabase: 0,
		Concurrency:       defaultConcurrency,
		PeriodMin:         defaultPeriodMin,
		PeriodMax:         defaultPeriodMax,
		Duration:          defaultDuration,
		Period:            defaultPeriod,
		TimeRange:         defaultTimeRange,
		HotPartitions:     defaultHotPartitions,
		TimeRangeMin:      defaultTimeRange,
		TimeRangeMax:      defaultTimeRangeMax,
		HotShare:          defaultHotShare,
		StatsInterval:     defaultStatsInterval,
		Profile:           defaultProfile,
		OutputFormat:      defaultOutputFormat,
		DryRun:            false,
		DiscoveryTimeout:  defaultDiscoveryTimeout,
		DialTimeout:       defaultDialTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
	}
}

// Validate checks that configuration values can produce safe, bounded queries.
func (c Config) Validate() error {
	if strings.TrimSpace(c.MySQLHost) == "" {
		return fmt.Errorf("mysql-host must not be empty")
	}
	if c.MySQLPort < 1 || c.MySQLPort > 65535 {
		return fmt.Errorf("mysql-port must be between 1 and 65535")
	}
	if len(c.Databases) == 0 || hasEmptyDatabase(c.Databases) {
		return fmt.Errorf("databases must contain at least one non-empty database")
	}
	if c.TablesPerDatabase < 0 {
		return fmt.Errorf("tables-per-database must be zero or greater")
	}
	if c.Concurrency < 1 {
		return fmt.Errorf("concurrency must be greater than zero")
	}
	if c.Duration < 0 || c.PeriodMin <= 0 || c.PeriodMax < c.PeriodMin || c.TimeRangeMin <= 0 || c.TimeRangeMax < c.TimeRangeMin || c.StatsInterval <= 0 {
		return fmt.Errorf("duration must be zero or greater, ranges must be positive with min <= max, and stats-interval must be greater than zero")
	}
	if c.DiscoveryTimeout <= 0 || c.DialTimeout <= 0 || c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
		return fmt.Errorf("discovery and MySQL timeouts must be greater than zero")
	}
	if c.Jitter < 0 {
		return fmt.Errorf("jitter must not be negative")
	}
	if c.HotPartitions < 1 {
		return fmt.Errorf("hot-partitions must be greater than zero")
	}
	if c.HotShare <= 0 || c.HotShare > 1 {
		return fmt.Errorf("hot-share must be greater than zero and at most one")
	}
	if c.Profile != "sustained" && c.Profile != "periodic" {
		return fmt.Errorf("profile must be sustained or periodic")
	}
	if c.OutputFormat != "ndjson" && c.OutputFormat != "csv" {
		return fmt.Errorf("output-format must be ndjson or csv")
	}
	return nil
}

// DSN returns the MySQL driver data source name without selecting a database.
func (c Config) DSN() string {
	dsnConfig := mysql.NewConfig()
	dsnConfig.User = c.MySQLUser
	dsnConfig.Passwd = c.MySQLPassword
	dsnConfig.Net = "tcp"
	dsnConfig.Addr = net.JoinHostPort(c.MySQLHost, strconv.Itoa(c.MySQLPort))
	dsnConfig.ParseTime = true
	dsnConfig.Timeout = c.DialTimeout
	dsnConfig.ReadTimeout = c.ReadTimeout
	dsnConfig.WriteTimeout = c.WriteTimeout
	return dsnConfig.FormatDSN()
}

// OpenDB creates a database/sql handle using the configured MySQL connection.
func (c Config) OpenDB() (*sql.DB, error) {
	return sql.Open("mysql", c.DSN())
}

func hasEmptyDatabase(databases []string) bool {
	for _, database := range databases {
		if strings.TrimSpace(database) == "" {
			return true
		}
	}
	return false
}

func parseDatabases(value string) []string {
	items := strings.Split(value, ",")
	databases := make([]string, 0, len(items))
	for _, item := range items {
		databases = append(databases, strings.TrimSpace(item))
	}
	return databases
}

type workloadFileConfig struct {
	Concurrency   *int     `yaml:"concurrency"`
	Duration      *string  `yaml:"duration"`
	Period        *string  `yaml:"period"`
	PeriodMin     *string  `yaml:"period_min"`
	PeriodMax     *string  `yaml:"period_max"`
	Jitter        *string  `yaml:"jitter"`
	TimeRange     *string  `yaml:"time_range"`
	TimeRangeMin  *string  `yaml:"time_range_min"`
	TimeRangeMax  *string  `yaml:"time_range_max"`
	HotPartitions *int     `yaml:"hot_partitions"`
	HotShare      *float64 `yaml:"hot_share"`
	StatsInterval *string  `yaml:"stats_interval"`
	Profile       *string  `yaml:"profile"`
}

func applyWorkloadFile(path string, config *Config, isFlagSet func(string) bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read workload config: %w", err)
	}
	var file workloadFileConfig
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&file); err != nil {
		return fmt.Errorf("parse workload config: %w", err)
	}
	if file.Concurrency != nil && !isFlagSet("concurrency") {
		config.Concurrency = *file.Concurrency
	}
	if file.HotPartitions != nil && !isFlagSet("hot-partitions") {
		config.HotPartitions = *file.HotPartitions
	}
	if file.HotShare != nil && !isFlagSet("hot-share") {
		config.HotShare = *file.HotShare
	}
	if file.Profile != nil && !isFlagSet("profile") {
		config.Profile = *file.Profile
	}
	for _, item := range []struct {
		name        string
		value       *string
		destination *time.Duration
	}{
		{"duration", file.Duration, &config.Duration}, {"period", file.Period, &config.Period}, {"period-min", file.PeriodMin, &config.PeriodMin}, {"period-max", file.PeriodMax, &config.PeriodMax}, {"jitter", file.Jitter, &config.Jitter}, {"time-range", file.TimeRange, &config.TimeRange}, {"time-range-min", file.TimeRangeMin, &config.TimeRangeMin}, {"time-range-max", file.TimeRangeMax, &config.TimeRangeMax}, {"stats-interval", file.StatsInterval, &config.StatsInterval},
	} {
		if item.value == nil || isFlagSet(item.name) {
			continue
		}
		value, err := time.ParseDuration(*item.value)
		if err != nil {
			return fmt.Errorf("parse workload config %s: %w", item.name, err)
		}
		*item.destination = value
	}
	if file.Period != nil && !isFlagSet("period-min") && !isFlagSet("period-max") {
		config.PeriodMin, config.PeriodMax = config.Period, config.Period
	}
	if file.TimeRange != nil && !isFlagSet("time-range-min") && !isFlagSet("time-range-max") {
		config.TimeRangeMin, config.TimeRangeMax = config.TimeRange, config.TimeRange
	}
	return nil
}
