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

	err := inputProcessor.writeFile(dataOutputPair, "userdata")

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

func TestProcessBase64GzippedCloudInitConfig(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)
	config := Config{outputDir: tempDir}
	inputProcessor := InputProcessor{config: config}
	dataOutputPairs := []DataOutputPair{
		{
			Data:      []byte("H4sICKnBLGUAA3VzZXJkYXRhAL2USY+jOBTH73yKqO6ptiFJhyrlEBYTEiCNAbPczFJAYpYJZIFP33R1qTQ96kOPRhpfbD37Lf+nn5/c1H1W93N3aLOXWXVlfdnSS/+lKh9Z+jqLm2ud0suwedr8ur5CuFiJS34F1gsI1wL/db3ZPHGmbqpzkl26sqlfZvAZcNx8/ieunPxLIX326L8krLmm86Sp38r8dZYU9NJl/ebp2r/N179L9RniQuvuLbvM1Tpp0rLOX2Yx7bLV4vOFUnZt05X9uyvte5oU1WR/nb2VLKtplW2e/p78eaAVe+I4fZBOsSbCSIanmF9W1E+umY/6eNsKScVARMTJ9mDJ2Bx0eXvSUTrQADN9h5vIkT58UphUCHDUF6+6FrWx5uWpts4NnoBwEMfQh/dYI29hIN2Ne5PrMsgTDQGqTOdSZKlm3hIe9Yn2YAGP7oksjlzMQ2YErIj9+8+KZP1HBbmukWvIi2fqL+tjKdXZpNk62R931i2uMYtre6WrOUgAU2wYIQ6fU2J76ujWUUkgmWSlWqrtFcslERnRaHve4KkQG8Da22dLOzpr3gGmkAAMsLAH1jm9myq8Ynf74EzAHg4kB6yh0pfFMBakFvMiNfno22FkyjFoTzEoFhFrsR3sBWdMPXzuTyHBHt0VpX1KLyZoXYoI4QxIePNsqd4OPUzlzFO/f/g7Uh9R1E0VOonQhnH1OKZnZGLeHnx1j6iKZAvogwv2SqREQgg93tXCJRd5xZAFbeMTycV8obgC7mMSwUiLDIqYSxXr7o6FmzIVWKMOEwUPrmcusQfPrlLIeESSWxcg8VDLEd/aHUBROjsTuv7yGCvenbpShbV+MSFzjRHRvF14PwyicQCopTt0csaidYR86fFoEbK2O3rISJk5cE5dfLMJkmz7txyB6B0RBKbzO2+T/abLIkwEfUKjKJLK+0SE+8nI+kegXt9JRarlK13Wb2lgDYZgNWGwZ9NexIHUTQAXyW47OYt95EAQBQUwarxIP3Di/jVPAdlj+A+c1NS3NKxwlisxByHJGXPBQ9gnvuiGrnSbevmXS7wFVZk52TTbs5hZRXuf7CXbQ5Kct0PqL0+xb6/kcpsbjhRz09d41z81+xT691tYseskC0zyJglWozOw+S/D6DHvioyxLrmUbf+/T6Oun4byc1dMQyjyrSYetqUn4CKZ2h3J+h8Lm8857js3myh88gUAAA=="),
			OutputDir: "i-111",
		},
	}
	err := inputProcessor.Process(dataOutputPairs)

	assert.FileExistsf(t, filepath.Join(tempDir, "i-111", "userdata"), "Expected userdata")
	assert.Nil(t, err)
}
