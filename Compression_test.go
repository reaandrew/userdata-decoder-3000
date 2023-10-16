package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsGzipped(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	assert.True(t, isGzipped([]byte(encryptedHelloWorld)))
}

func TestUnzip(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld

	result, err := unzipData([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}
