package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsGzipped(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	assert.True(t, isGzipped([]byte(encryptedHelloWorld)))
}

func TestIsNotGzipped(t *testing.T) {
	assert.False(t, isGzipped([]byte("Data")))
}

func TestUnzip(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld

	result, err := unzipData([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestUnzipReturnsError(t *testing.T) {
	_, err := unzipData([]byte("Value"))
	assert.NotNil(t, err)
}
