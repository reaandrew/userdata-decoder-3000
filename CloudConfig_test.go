package main

import (
	_ "bytes"
	"os"
	"path/filepath"
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

func TestSaveWriteFiles(t *testing.T) {
	tests := []struct {
		name          string
		cloudConfig   CloudConfig
		expectedFiles map[string]string
		expectedError error
	}{
		{
			name: "Save multiple files",
			cloudConfig: CloudConfig{
				WriteFiles: []WriteFile{
					{Path: "file1.txt", Content: "content1"},
					{Path: "dir/file2.txt", Content: "content2"},
				},
			},
			expectedFiles: map[string]string{
				"file1.txt":     "content1",
				"dir/file2.txt": "content2",
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("could not create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir) // clean up

			// Run the function
			err = tt.cloudConfig.SaveWriteFiles(tempDir)

			// Check for expected error
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}

			// Check for unexpected error
			assert.NoError(t, err)

			// Validate the files were saved correctly
			for relPath, expectedContent := range tt.expectedFiles {
				fullPath := filepath.Join(tempDir, relPath)
				actualContent, err := os.ReadFile(fullPath)
				if err != nil {
					t.Fatalf("could not read file: %v", err)
				}
				assert.Equal(t, expectedContent, string(actualContent))
			}
		})
	}
}
