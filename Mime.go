package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"strings"
)

var errFailedToExtractMimeBoundary = errors.New("failed to get MIME boundary")

type MimeAttachment struct {
	ContentType string
	Filename    string
	Content     []byte
}

func decodeMimAttachments(data []byte) (attachments []MimeAttachment, err error) {
	boundary, err := extractBoundary(data)
	if err != nil {
		return nil, errFailedToExtractMimeBoundary
	}

	reader := multipart.NewReader(bytes.NewReader(data), boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break // End of the stream, no more parts
			}
			return nil, fmt.Errorf("error reading part: %s", err) // Actual error reading part
		}

		contentType := part.Header.Get("Content-Type")
		disposition, _, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		if err != nil {
			return nil, fmt.Errorf("error parsing content disposition: %s", err)
		}

		content, err := io.ReadAll(part)
		if err != nil {
			return nil, fmt.Errorf("error reading part content: %s", err)
		}

		filename := ""
		if disposition == "attachment" {
			_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
			if err == nil {
				filename = params["filename"]
			}
		}

		attachments = append(attachments, MimeAttachment{
			ContentType: contentType,
			Filename:    filename,
			Content:     content,
		})

		// Remember to close each part after reading
		part.Close()
	}
	return attachments, nil
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

func ExtractMimeAttachmentsFromBytes(decoded []byte) (attachments []MimeAttachment, err error) {
	return decodeMimAttachments(decoded)
}
