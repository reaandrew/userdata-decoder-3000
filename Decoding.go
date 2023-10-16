package main

import "encoding/base64"

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

func decode(data []byte) ([]byte, error) {
	if isGzipped(data) {
		return unzipData(data)
	}
	return decodeBase64(data)
}
