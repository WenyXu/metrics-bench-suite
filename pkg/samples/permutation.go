package samples

import (
	"fmt"
	"maps"
)

// Label is a struct that represents a label in a metric.
type Label struct {
	Name   string
	Values []string
}

// TagSetPermutation generates all possible permutations of tag sets.
func TagSetPermutation(labels []Label, totalCount *int) []map[string]string {
	Set := make(map[string][]string)
	for _, label := range labels {
		Set[label.Name] = label.Values
	}

	keys := make([]string, 0, len(Set))
	values := make([][]string, 0, len(Set))
	for k, v := range Set {
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
