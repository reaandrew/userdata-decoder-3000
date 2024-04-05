package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractMimeAttachments(t *testing.T) {
	decoded, _ := decode([]byte(mimeMessage))

	attachments, err := ExtractMimeAttachmentsFromBytes(decoded)

	assert.Nil(t, err)
	assert.Len(t, attachments, 2)
	assert.Equal(t, attachments[0].ContentType, "text/cloud-config; charset=\"utf-8\"")
	assert.Equal(t, attachments[1].ContentType, "text/x-shellscript; charset=\"utf-8\"")
}

func TestExtractMimeFilenames(t *testing.T) {
	decoded, _ := decode([]byte(mimeMessage))

	attachments, err := ExtractMimeAttachmentsFromBytes(decoded)

	assert.Nil(t, err)
	assert.Len(t, attachments, 2)
	assert.Equal(t, attachments[0].Filename, "cloud-config.yaml")
	assert.Equal(t, attachments[1].Filename, "start.sh")
}

func TestExtractMimeAttachmentsReturnsErrorWhenNotMime(t *testing.T) {
	_, err := ExtractMimeAttachmentsFromBytes([]byte("not mime"))

	assert.NotNil(t, err)
	assert.Equal(t, err, errFailedToExtractMimeBoundary)
}
