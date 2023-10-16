package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeBase64(t *testing.T) {
	result, err := decodeBase64([]byte(base64EncodedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestDecodeWithEncrypted(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestDecode(t *testing.T) {
	encryptedHelloWorld := base64EncodedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestDecodeNonEncodedValue(t *testing.T) {
	expected := "Hello, World!"
	result, err := decode([]byte(expected))
	assert.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
