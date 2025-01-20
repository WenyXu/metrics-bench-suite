package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var base64Data = "U/BSClEKHwoIX19uYW1lX18SE2h0dHBfcmVxdWVzdHNfdG90YWwKDQoGc3RhdHVzEgMyMDAKDQoGbWV0aG9kEgNHRVQSEAkAAAAAAADwPxCMuO6ZyDI="

func TestDecodeBody(t *testing.T) {
	compressed, err := Base64(base64Data)
	assert.Nil(t, err)
	_, err = Body(compressed)
	assert.Nil(t, err)
}
