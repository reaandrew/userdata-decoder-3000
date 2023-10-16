package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path"
	"strings"
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

type MimeAttachment struct {
	ContentType string
	Content     []byte
}

func extractBoundary(data []byte) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Content-Type:") {
			_, params, err := mime.ParseMediaType(strings.TrimPrefix(line, "Content-Type: "))
			if err != nil {
				return "", err
			}
			boundary, found := params["boundary"]
			if !found {
				return "", io.EOF
			}
			return boundary, nil
		}
	}
	return "", io.EOF
}

func decodeMimAttachments(data []byte) (attachments []MimeAttachment, err error) {
	boundary, err := extractBoundary(data)
	if err != nil {
		log.Fatalf("Failed to get boundary: %s", err)
	}

	reader := multipart.NewReader(bytes.NewReader(data), boundary)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading part: ", err)
		}

		contentType := part.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		content, err := io.ReadAll(part)
		if err != nil {
			log.Fatal("Error reading content: ", err)
		}

		encoding := part.Header.Get("Content-Transfer-Encoding")
		if encoding == "base64" {
			decoded, err := base64.StdEncoding.DecodeString(string(content))
			if err != nil {
				log.Fatal("Error decoding base64: ", err)
			}
			content = decoded
		}

		attachment := MimeAttachment{
			ContentType: contentType,
			Content:     content,
		}

		attachments = append(attachments, attachment)
	}

	return attachments, err
}

func extractMimeAttachments(encodedData []byte) (attachments []MimeAttachment, err error) {
	decoded, err := decode(encodedData)
	if err != nil {

		return
	}

	return decodeMimAttachments(decoded)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	output, input, err := parseFlags()
	if err != nil {
		return err
	}

	attachments, err := extractMimeAttachments([]byte(input))
	if err != nil {
		return fmt.Errorf("failed to extract mime attachments: %w", err)
	}

	for _, attachment := range attachments {
		if strings.Contains(attachment.ContentType, "text/cloud-config") {
			cloudConfig, err := ReadCloudConfigFrom(attachment, output)
			if err != nil {
				return fmt.Errorf("failed to extract cloud config write files: %w", err)
			}

			if err := saveWriteFiles(cloudConfig.WriteFiles); err != nil {
				return fmt.Errorf("failed to save write files: %w", err)
			}
		}
	}

	return nil
}

func parseFlags() (output string, input string, err error) {
	flag.StringVar(&output, "o", "output", "Specify the output directory within your working directory.")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		return "", "", errors.New("input file is required")
	}

	return output, args[0], nil
}

func saveWriteFiles(writeFiles []WriteFile) error {
	for _, file := range writeFiles {
		err := os.MkdirAll(path.Dir(file.Path), 0755)
		if err != nil {
			return fmt.Errorf("error creating output directories: %w", err)
		}
		err = os.WriteFile(file.Path, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
	}
	return nil
}
