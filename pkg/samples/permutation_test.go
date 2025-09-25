package samples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagSetPermutation(t *testing.T) {
	metrics := []LabelCandidates{
		{Name: "label1", Values: []string{"value1", "value2"}},
		{Name: "label2", Values: []string{"value3", "value4"}},
		{Name: "label3", Values: []string{"value5", "value6"}},
	}

	total := 0
	permutations := TagSetPermutation(metrics, &total)
	assert.Equal(t, len(permutations), 8)
	assert.Equal(t, total, 8)
}
