package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandLineProviderFetchData(t *testing.T) {
	provider := CommandLineProvider{
		Input: "test_input",
	}

	outputPairs, err := provider.FetchData()

	assert.NoError(t, err)
	assert.NotNil(t, outputPairs)
	assert.Equal(t, 1, len(outputPairs))
	assert.Equal(t, "test_input", string(outputPairs[0].Data))
	assert.Equal(t, "", outputPairs[0].OutputDir)
}
