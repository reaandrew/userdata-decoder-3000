package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
)

func decodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		encoded := scanner.Text()
		decoded, err := decodeBase64(encoded)
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
