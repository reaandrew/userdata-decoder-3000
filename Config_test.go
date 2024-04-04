package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlesUnknownProvider(t *testing.T) {
	_, err := Config{providerKey: "talula"}.getProvider()
	assert.ErrorContains(t, err, "unknown provider")
}

func TestReturnsCommandLineProvider(t *testing.T) {
	provider, err := Config{
		outputDir:   "",
		providerKey: "",
		args: []string{
			"something",
		},
	}.getProvider()

	assert.Nil(t, err)
	assert.IsType(t, CommandLineProvider{}, provider)
}

func TestReturnsAwsProvider(t *testing.T) {
	provider, err := Config{
		outputDir:   "",
		providerKey: "aws",
	}.getProvider()

	assert.Nil(t, err)
	assert.IsType(t, &AWSProvider{}, provider)
}
