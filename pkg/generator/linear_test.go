package generator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinearTrend(t *testing.T) {
	coef := 0.05
	offset := 1.0

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	lt, err := NewLinearTrend(coef, offset, startDate, time.Hour*24)
	assert.NoError(t, err)

	data, err := lt.Generate()
	assert.NoError(t, err)

	assert.Equal(t, data.Timestamp, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, data.Value, 2.05)

	data, err = lt.Generate()
	assert.NoError(t, err)

	assert.Equal(t, data.Timestamp, time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, data.Value, 2.1)

	data, err = lt.Generate()
	assert.NoError(t, err)

	assert.Equal(t, data.Timestamp, time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, data.Value, 2.15)

	data, err = lt.Generate()
	assert.NoError(t, err)

	assert.Equal(t, data.Timestamp, time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, data.Value, 2.2)

	data, err = lt.Generate()
	assert.NoError(t, err)

	assert.Equal(t, data.Timestamp, time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, data.Value, 2.25)
}

func TestLinearTrend2(t *testing.T) {
	coef := 0.05
	offset := 1.0

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	lt, err := NewLinearTrend(coef, offset, startDate, time.Second*10)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		data, err := lt.Generate()
		assert.NoError(t, err)
		fmt.Println(data)
	}
}
