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

func extractFilename(part *multipart.Part) (string, error) {
	// Retrieve the Content-Disposition header of the current part
	cdHeader := part.Header.Get("Content-Disposition")
	if cdHeader == "" {
		return "", fmt.Errorf("no Content-Disposition header found")
	}

	// Parse the Content-Disposition header to get disposition type and parameters
	_, params, err := mime.ParseMediaType(cdHeader)
	if err != nil {
		return "", fmt.Errorf("failed to parse Content-Disposition header: %w", err)
	}

	// The filename parameter should be present in the parameters map
	filename, found := params["filename"]
	if !found {
		return "", fmt.Errorf("filename not found in Content-Disposition header")
	}

	return filename, nil
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

		// Assuming extractContentType and io.ReadAll are defined or replaced appropriately
		contentType := part.Header.Get("Content-Type")

		content, err := io.ReadAll(part)
		if err != nil {
			return nil, fmt.Errorf("error reading part content: %s", err)
		}

		//filename, err := extractFilename(part)
		//if err != nil {
		//	return nil, fmt.Errorf("error reading part filename: %s", err)
		//}

		attachments = append(attachments, MimeAttachment{
			ContentType: contentType,
			//Filename:    filename,
			Content: content,
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
