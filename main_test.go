package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DecodeBase64(t *testing.T) {
	result, err := decodeBase64("SGVsbG8sIFdvcmxkIQ==")
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}
