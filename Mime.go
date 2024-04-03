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
	Content     []byte
}

func decodeMimAttachments(data []byte) (attachments []MimeAttachment, err error) {
	mimeData, err := decode(data)
	boundary, err := extractBoundary(mimeData)
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

		attachments = append(attachments, MimeAttachment{
			ContentType: contentType,
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
	//boundary, err := extractBoundary(decoded)
	//if err != nil {
	//	return []MimeAttachment{}, errFailedToExtractMimeBoundary
	//}
	//
	//reader := multipart.NewReader(bytes.NewReader(decoded), boundary)
	//
	//for {
	//	part, err := reader.NextPart()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		log.Fatal("Error reading part: ", err)
	//	}
	//
	//	contentType := part.Header.Get("Content-Type")
	//	if contentType == "" {
	//		contentType = "application/octet-stream"
	//	}
	//
	//	content, err := io.ReadAll(part)
	//	if err != nil {
	//		log.Fatal("Error reading content: ", err)
	//	}
	//
	//	encoding := part.Header.Get("Content-Transfer-Encoding")
	//	if encoding == "base64" {
	//		decoded, err := base64.StdEncoding.DecodeString(string(content))
	//		if err != nil {
	//			log.Fatal("Error decoding base64: ", err)
	//		}
	//		content = decoded
	//	}
	//
	//	attachment := MimeAttachment{
	//		ContentType: contentType,
	//		Content:     content,
	//	}
	//
	//	attachments = append(attachments, attachment)
	//}
	//
	//return attachments, err
	return decodeMimAttachments(decoded)
}
