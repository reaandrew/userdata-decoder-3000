package main

import (
	_ "bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCloudConfigFrom(t *testing.T) {
	tests := []struct {
		name          string
		attachment    MimeAttachment
		baseDir       string
		expected      CloudConfig
		expectedError string
	}{
		{
			name: "Valid cloud-config",
			attachment: MimeAttachment{
				ContentType: "text/cloud-config",
				Content:     []byte("write_files:\n  - path: /test.txt\n    encoding: b64\n    content: dGVzdA=="),
			},
			baseDir: "/base",
			expected: CloudConfig{
				WriteFiles: []WriteFile{
					{
						Path:     "/test.txt",
						Encoding: "b64",
						Content:  "test",
					},
				},
			},
		},
		{
			name: "Invalid content type",
			attachment: MimeAttachment{
				ContentType: "text/plain",
				Content:     []byte("write_files:\n  - path: /test.txt\n    encoding: b64\n    content: dGVzdA=="),
			},
			baseDir:       "/base",
			expectedError: "not a cloud-config content type",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReadCloudConfigFrom(tt.attachment)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
