package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"
)

var errFailedToExtractMimeBoundary = errors.New("failed to get MIME boundary")

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

func ExtractMimeAttachmentsFromBytes(data []byte) (attachments []MimeAttachment, err error) {
	boundary, err := extractBoundary(data)
	if err != nil {
		return []MimeAttachment{}, errFailedToExtractMimeBoundary
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
