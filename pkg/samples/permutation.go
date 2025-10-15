package samples

import (
	"fmt"
	"maps"
)

// LabelCandidates is a struct that represents a label in a metric.
type LabelCandidates struct {
	Name   string
	Values []string
}

// TagSetPermutation generates all possible permutations of tag sets.
func TagSetPermutation(labels []LabelCandidates, totalCount *int) []map[string]string {
	// Label name -> candidates.
	candidatesMap := make(map[string][]string)
	for _, label := range labels {
		candidatesMap[label.Name] = label.Values
	}

	keys := make([]string, 0, len(candidatesMap))
	values := make([][]string, 0, len(candidatesMap))
	for k, v := range candidatesMap {
		keys = append(keys, k)
		values = append(values, v)
	}

	count := 1
	for _, v := range values {
		count *= len(v)
	}
	fmt.Println("number of tag combinations:", count)
	*totalCount += count

	var permutations []map[string]string
	generatePermutations(keys, values, map[string]string{}, &permutations, 0)
	return permutations
}

// generatePermutations is a recursive helper function to generate permutations.
func generatePermutations(keys []string, values [][]string, current map[string]string, permutations *[]map[string]string, depth int) {
	if depth == len(keys) {
		// Make a copy of the current map and add it to the permutations
		perm := make(map[string]string)
		maps.Copy(perm, current)
		*permutations = append(*permutations, perm)
		return
	}

	for _, value := range values[depth] {
		current[keys[depth]] = value
		generatePermutations(keys, values, current, permutations, depth+1)
	}
}
