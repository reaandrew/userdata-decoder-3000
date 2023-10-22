package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandLineProviderFetchData(t *testing.T) {
	// Arrange
	provider := CommandLineProvider{
		Input: "test_input",
	}

	// Act
	outputPairs, err := provider.FetchData()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, outputPairs)
	assert.Equal(t, 1, len(outputPairs))
	assert.Equal(t, "test_input", string(outputPairs[0].Data))
	assert.Equal(t, "", outputPairs[0].OutputDir)
}
