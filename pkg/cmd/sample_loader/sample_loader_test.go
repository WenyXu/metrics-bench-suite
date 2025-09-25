package sample_loader

import (
	"metrics-bench-suite/pkg/samples"
	"reflect"
	"sort"
	"testing"
)

func TestTagSetPermutationStream(t *testing.T) {
	// Test case 1: Empty labels
	t.Run("Empty labels", func(t *testing.T) {
		labels := []samples.LabelCandidates{}
		permChan := make(chan map[string]string, 10)
		totalCount := 0

		go TagSetPermutationStream(labels, permChan, &totalCount)

		println("1")
		results := make([]map[string]string, 0)
		for perm := range permChan {
			println("2")
			results = append(results, perm)
		}

		println("3")
		expectedCount := 1
		if totalCount != expectedCount {
			t.Errorf("Expected total count %d, got %d", expectedCount, totalCount)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		expected := map[string]string{}
		if !reflect.DeepEqual(results[0], expected) {
			t.Errorf("Expected %v, got %v", expected, results[0])
		}
	})

	// Test case 2: Single label with one value
	t.Run("Single label with one value", func(t *testing.T) {
		labels := []samples.LabelCandidates{
			{
				Name:   "label1",
				Values: []string{"value1"},
			},
		}
		permChan := make(chan map[string]string, 10)
		totalCount := 0

		go TagSetPermutationStream(labels, permChan, &totalCount)

		results := make([]map[string]string, 0)
		for perm := range permChan {
			results = append(results, perm)
		}

		expectedCount := 1
		if totalCount != expectedCount {
			t.Errorf("Expected total count %d, got %d", expectedCount, totalCount)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		expected := map[string]string{"label1": "value1"}
		if !reflect.DeepEqual(results[0], expected) {
			t.Errorf("Expected %v, got %v", expected, results[0])
		}
	})

	// Test case 3: Single label with multiple values
	t.Run("Single label with multiple values", func(t *testing.T) {
		labels := []samples.LabelCandidates{
			{
				Name:   "label1",
				Values: []string{"value1", "value2", "value3"},
			},
		}
		permChan := make(chan map[string]string, 10)
		totalCount := 0

		go TagSetPermutationStream(labels, permChan, &totalCount)

		results := make([]map[string]string, 0)
		for perm := range permChan {
			results = append(results, perm)
		}

		expectedCount := 3
		if totalCount != expectedCount {
			t.Errorf("Expected total count %d, got %d", expectedCount, totalCount)
		}

		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}

		// Sort results for consistent comparison
		sort.Slice(results, func(i, j int) bool {
			return results[i]["label1"] < results[j]["label1"]
		})

		expected := []map[string]string{
			{"label1": "value1"},
			{"label1": "value2"},
			{"label1": "value3"},
		}
		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	// Test case 4: Multiple labels
	t.Run("Multiple labels", func(t *testing.T) {
		labels := []samples.LabelCandidates{
			{
				Name:   "label1",
				Values: []string{"a", "b"},
			},
			{
				Name:   "label2",
				Values: []string{"x", "y"},
			},
		}
		permChan := make(chan map[string]string, 10)
		totalCount := 0

		go TagSetPermutationStream(labels, permChan, &totalCount)

		results := make([]map[string]string, 0)
		for perm := range permChan {
			results = append(results, perm)
		}

		expectedCount := 4 // 2 * 2 = 4 combinations
		if totalCount != expectedCount {
			t.Errorf("Expected total count %d, got %d", expectedCount, totalCount)
		}

		if len(results) != 4 {
			t.Errorf("Expected 4 results, got %d", len(results))
		}

		// Sort results for consistent comparison
		sort.Slice(results, func(i, j int) bool {
			if results[i]["label1"] != results[j]["label1"] {
				return results[i]["label1"] < results[j]["label1"]
			}
			return results[i]["label2"] < results[j]["label2"]
		})

		expected := []map[string]string{
			{"label1": "a", "label2": "x"},
			{"label1": "a", "label2": "y"},
			{"label1": "b", "label2": "x"},
			{"label1": "b", "label2": "y"},
		}
		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})

	// Test case 5: Multiple labels with different value counts
	t.Run("Multiple labels with different value counts", func(t *testing.T) {
		labels := []samples.LabelCandidates{
			{
				Name:   "env",
				Values: []string{"prod", "dev"},
			},
			{
				Name:   "region",
				Values: []string{"us-east", "us-west", "eu-west"},
			},
			{
				Name:   "service",
				Values: []string{"web"},
			},
		}
		permChan := make(chan map[string]string, 100) // Larger buffer to handle all combinations
		totalCount := 0

		go TagSetPermutationStream(labels, permChan, &totalCount)

		results := make([]map[string]string, 0)
		for perm := range permChan {
			results = append(results, perm)
		}

		expectedCount := 6 // 2 * 3 * 1 = 6 combinations
		if totalCount != expectedCount {
			t.Errorf("Expected total count %d, got %d", expectedCount, totalCount)
		}

		if len(results) != 6 {
			t.Errorf("Expected 6 results, got %d", len(results))
		}

		// Sort results for consistent comparison
		sort.Slice(results, func(i, j int) bool {
			if results[i]["env"] != results[j]["env"] {
				return results[i]["env"] < results[j]["env"]
			}
			if results[i]["region"] != results[j]["region"] {
				return results[i]["region"] < results[j]["region"]
			}
			return results[i]["service"] < results[j]["service"]
		})

		expected := []map[string]string{
			{"env": "dev", "region": "eu-west", "service": "web"},
			{"env": "dev", "region": "us-east", "service": "web"},
			{"env": "dev", "region": "us-west", "service": "web"},
			{"env": "prod", "region": "eu-west", "service": "web"},
			{"env": "prod", "region": "us-east", "service": "web"},
			{"env": "prod", "region": "us-west", "service": "web"},
		}
		if !reflect.DeepEqual(results, expected) {
			t.Errorf("Expected %v, got %v", expected, results)
		}
	})
}
