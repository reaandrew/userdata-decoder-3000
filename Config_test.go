package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlesUnknownProvider(t *testing.T) {
	_, err := Config{providerKey: "talula"}.getProvider()
	assert.ErrorContains(t, err, "unknown provider")
}
