package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func decodeBase64(data []byte) ([]byte, error) {
	var output []byte
	maxDecodedLen := base64.StdEncoding.DecodedLen(len(data))
	output = make([]byte, maxDecodedLen)
	n, err := base64.StdEncoding.Decode(output, data)
	if err != nil {
		return nil, err
	}

	return output[:n], nil
}

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
		return nil, err
	}
	return io.ReadAll(r)
}

func decode(data []byte) ([]byte, error) {
	if isGzipped(data) {
		return unzipData(data)
	}
	return decodeBase64(data)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		encoded := scanner.Text()
		decoded, err := decode([]byte(encoded))
		if err != nil {
			fmt.Println("Failed to decode base64:", err)
			return
		}
		fmt.Printf("%s", string(decoded))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}
