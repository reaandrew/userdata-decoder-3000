package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
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

type CloudConfig struct {
	WriteFiles []WriteFile `yaml:"write_files"`
}

type WriteFile struct {
	Path     string `yaml:"path"`
	Encoding string `yaml:"encoding"`
	Content  string `yaml:"content"`
}

func extractCloudConfigWriteFiles(attachment MimeAttachment, baseDir string) (writefiles []WriteFile, err error) {
	if !strings.Contains(attachment.ContentType, "text/cloud-config") {
		return []WriteFile{}, fmt.Errorf("not a cloud-config content type")
	}

	var config CloudConfig

	err = yaml.Unmarshal(attachment.Content, &config)
	if err != nil {
		return []WriteFile{}, err
	}

	for _, file := range config.WriteFiles {
		fullPath := filepath.Join(baseDir, file.Path)

		dir := filepath.Dir(fullPath)
		content, err := decode([]byte(file.Content))
		if err != nil {
			return []WriteFile{}, err
		}

		writefiles = append(writefiles, WriteFile{
			Path:    dir,
			Content: string(content),
		})
	}
	return writefiles, err
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		encoded := scanner.Text()
		attachments, err := extractMimeAttachments([]byte(encoded))

		if err != nil {
			fmt.Println("Failed to extract mime attachments")
		}

		for _, attachment := range attachments {
			fmt.Printf("Attachment - Media Type: %s \n", attachment.ContentType)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}
