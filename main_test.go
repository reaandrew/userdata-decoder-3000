package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	base64EncodedHelloWorld       = "SGVsbG8sIFdvcmxkIQ=="
	base64EncodedZippedHelloWorld = "H4sIAAAAAAAAA/NIzcnJ11EIzy/KSVEEANDDSuwNAAAA"
)

func Test_DecodeBase64(t *testing.T) {
	result, err := decodeBase64([]byte(base64EncodedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}

func Test_IsGzipped(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	assert.True(t, isGzipped([]byte(encryptedHelloWorld)))
}

func Test_Unzip(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld

	result, err := unzipData([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}

func Test_Decode_With_Encrypted(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}

func Test_Decode(t *testing.T) {
	encryptedHelloWorld := base64EncodedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}
