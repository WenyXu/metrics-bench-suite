package partition_query_loader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultConfigIsValid(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("expected defaults to be valid, got %v", err)
	}
}

func TestDefaultConfigRunsUntilCancelled(t *testing.T) {
	config := DefaultConfig()
	if config.DryRun || config.Duration != 0 {
		t.Fatalf("defaults dry_run=%t duration=%s", config.DryRun, config.Duration)
	}
}

func TestDefaultConfigUsesRebalanceFriendlySkewRanges(t *testing.T) {
	config := DefaultConfig()
	if config.PeriodMin != 15*time.Second || config.PeriodMax != 45*time.Second || config.TimeRangeMin != 5*time.Minute || config.TimeRangeMax != 15*time.Minute || config.HotPartitions != 3 || config.HotShare != 0.85 {
		t.Fatalf("unexpected skew defaults: %#v", config)
	}
}

func TestWorkloadConfigAppliesYAMLAndFlagsOverrideIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workload.yaml")
	if err := os.WriteFile(path, []byte("concurrency: 3\nduration: 20s\nperiod: 2s\nprofile: periodic\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var got Config
	command := NewCommandWithRunner(func(config Config) error { got = config; return nil })
	command.SetArgs([]string{"--workload-config", path, "--period", "1s"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if got.Concurrency != 3 || got.Duration != 20*time.Second || got.Period != time.Second || got.Profile != "periodic" {
		t.Fatalf("unexpected workload config: %#v", got)
	}
}

func TestConfigValidateRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{"empty mysql host", func(config *Config) { config.MySQLHost = " " }},
		{"invalid mysql port", func(config *Config) { config.MySQLPort = 0 }},
		{"empty databases", func(config *Config) { config.Databases = nil }},
		{"negative table limit", func(config *Config) { config.TablesPerDatabase = -1 }},
		{"zero concurrency", func(config *Config) { config.Concurrency = 0 }},
		{"negative duration", func(config *Config) { config.Duration = -time.Second }},
		{"negative jitter", func(config *Config) { config.Jitter = -time.Second }},
		{"zero hot partitions", func(config *Config) { config.HotPartitions = 0 }},
		{"invalid hot share", func(config *Config) { config.HotShare = 1.1 }},
		{"unknown profile", func(config *Config) { config.Profile = "burst" }},
		{"unknown output format", func(config *Config) { config.OutputFormat = "json" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := DefaultConfig()
			test.mutate(&config)
			if err := config.Validate(); err == nil {
				t.Fatal("expected invalid configuration error")
			}
		})
	}
}

func TestConfigFromCommandParsesFlags(t *testing.T) {
	var got Config
	command := NewCommandWithRunner(func(config Config) error {
		got = config
		return nil
	})
	command.SetArgs([]string{
		"--mysql-host", "greptime.example",
		"--mysql-port", "4100",
		"--databases", "metrics1, metrics2",
		"--tables-per-database", "7",
		"--random-seed", "42",
		"--concurrency", "4",
		"--duration", "2m",
		"--period", "10s",
		"--jitter", "2s",
		"--time-range", "30m",
		"--hot-partitions", "5",
		"--hot-share", "0.9",
		"--stats-interval", "20s",
		"--profile", "periodic",
		"--output-format", "csv",
		"--dry-run=false",
	})
	if err := command.Execute(); err != nil {
		t.Fatalf("expected valid configuration, got %v", err)
	}
	if got.MySQLHost != "greptime.example" || got.MySQLPort != 4100 {
		t.Fatalf("unexpected MySQL endpoint: %s:%d", got.MySQLHost, got.MySQLPort)
	}
	if len(got.Databases) != 2 || got.Databases[1] != "metrics2" {
		t.Fatalf("unexpected databases: %#v", got.Databases)
	}
	if got.Profile != "periodic" || got.OutputFormat != "csv" || got.DryRun {
		t.Fatalf("unexpected parsed options: %#v", got)
	}
	if got.TablesPerDatabase != 7 || got.RandomSeed == nil || *got.RandomSeed != 42 {
		t.Fatalf("unexpected discovery selection options: %#v", got)
	}
}

func TestConfigDSNUsesConfiguredEndpointWithoutDatabase(t *testing.T) {
	config := DefaultConfig()
	config.MySQLHost = "greptime.example"
	config.MySQLPort = 4100

	dsn := config.DSN()
	if !strings.Contains(dsn, "tcp(greptime.example:4100)/") {
		t.Fatalf("unexpected DSN: %q", dsn)
	}
	if !strings.Contains(dsn, "parseTime=true") {
		t.Fatalf("expected parseTime DSN parameter: %q", dsn)
	}
}

func TestCommandReturnsInvalidFlagConfiguration(t *testing.T) {
	command := NewCommand()
	command.SetArgs([]string{"--hot-share", "0"})
	if err := command.Execute(); err == nil {
		t.Fatal("expected invalid flag configuration error")
	}
}
