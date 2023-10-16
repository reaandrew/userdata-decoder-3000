package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func isGzipped(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	decodedData, err := decodeBase64(data)
	if err != nil {
		return false
	}
	return len(decodedData) > 1 && decodedData[0] == 0x1F && decodedData[1] == 0x8B
}

func unzipData(data []byte) ([]byte, error) {
	decodedData, err := decodeBase64(data)
	r, err := gzip.NewReader(bytes.NewBuffer(decodedData))
	if err != nil {
		return nil, fmt.Errorf("failed to unzip data: %w", err)
	}
	return io.ReadAll(r)
}
