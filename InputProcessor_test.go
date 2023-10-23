package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWritePlainUserDataFileWhenNotBase64Encoded(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)

	config := Config{outputDir: tempDir}
	inputProcessor := InputProcessor{config: config}

	dataOutputPair := DataOutputPair{
		Data:      []byte("some data"),
		OutputDir: "some_output_dir",
	}

	err := inputProcessor.writePlainUserDataFile(dataOutputPair)

	assert.Nil(t, err)

	expectedFilePath := filepath.Join(tempDir, "some_output_dir", "userdata")
	_, err = os.Stat(expectedFilePath)
	assert.False(t, os.IsNotExist(err))
}

func TestWritePlainUserDataFileWhenNotMime(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)

	config := Config{outputDir: tempDir}
	inputProcessor := InputProcessor{config: config}

	dataOutputPairs := []DataOutputPair{
		{
			Data:      []byte(base64EncodedHelloWorld),
			OutputDir: "some_output_dir1",
		},
	}

	err := inputProcessor.Process(dataOutputPairs)

	assert.Nil(t, err)
}

func TestProcess(t *testing.T) {

	tempDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)

	config := Config{outputDir: tempDir}
	inputProcessor := InputProcessor{config: config}

	dataOutputPairs := []DataOutputPair{
		{
			Data:      []byte("some data"),
			OutputDir: "some_output_dir1",
		},
		{
			Data:      []byte("some other data"),
			OutputDir: "some_output_dir2",
		},
	}

	err := inputProcessor.Process(dataOutputPairs)

	assert.Nil(t, err)
}
