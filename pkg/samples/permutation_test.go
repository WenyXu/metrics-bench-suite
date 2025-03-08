package samples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagSetPermutation(t *testing.T) {
	metrics := []Label{
		{Name: "label1", Values: []string{"value1", "value2"}},
		{Name: "label2", Values: []string{"value3", "value4"}},
		{Name: "label3", Values: []string{"value5", "value6"}},
	}

	permutations := TagSetPermutation(metrics)
	assert.Equal(t, len(permutations), 8)
}
