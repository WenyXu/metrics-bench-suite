package samples

import (
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents a parsed YAML configuration file
type FileConfig struct {
	Name   string
	Config Config
}

func getFileNameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return base[:len(base)-len(ext)]
}

func WalkAndParseConfigWithMaxFileCount(path string, tablePickCount uint64) ([]FileConfig, error) {
	var fileConfigs []FileConfig

	var totalSeries = 0
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			data, err := parseYAML(path)
			if err != nil {
				log.Printf("Error parsing YAML file %s: %v\n", path, err)
				return nil
			}

			metricName := getFileNameWithoutExt(path)
			num_series := 1
			for _, tag := range data.Tags {
				num_series *= tag.Dist.LabelGenerator().NumCandidates()
			}
			log.Printf("Parsing file: %s, num series: %d\n", path, num_series)
			totalSeries += num_series

			fileConfigs = append(fileConfigs, FileConfig{
				Name:   metricName,
				Config: data,
			})
			if uint64(len(fileConfigs)) > tablePickCount {
				log.Printf("Warning: More than %d YAML files found. Only the first %d will be used.\n", tablePickCount, tablePickCount)
				return fs.SkipAll
			}
		}
		return nil
	})

	log.Printf("Total series: %d\n", totalSeries)

	if err != nil {
		return nil, err
	}

	return fileConfigs, nil
}

// WalkAndParseConfig walks a directory and parses all YAML files, returning a list of FileConfig
func WalkAndParseConfig(path string) ([]FileConfig, error) {
	return WalkAndParseConfigWithMaxFileCount(path, math.MaxUint64)
}

// parseYAML parses a YAML file and returns a Config
func parseYAML(path string) (Config, error) {
	var config Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
