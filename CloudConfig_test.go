package main

import (
	_ "bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("could not create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			err = tt.cloudConfig.SaveWriteFiles(tempDir)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)

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

func TestExtractCloudConfig(t *testing.T) {
	tests := []struct {
		name          string
		attachments   []MimeAttachment
		expectedFiles map[string]string
		expectedError error
	}{
		{
			name: "Extract and save multiple files",
			attachments: []MimeAttachment{
				{
					ContentType: "text/cloud-config",
					Content: []byte(`---
write_files:
  - path: file1.txt
    content: H4sICAuIGWUAA3NvbWUtdGV4dC50eHQAK87PTVUoSa0oUcjMU8hILUoFACe2fIYRAAAA
    encoding: gz+b64
  - path: dir/file2.txt
    content: H4sICLDULGUAA3NvbWUtdGV4dC50eHQAK87PTVXIL8lILVIoSa0oUcjMUwCyUwECYhs6FwAAAA==
    encoding: gz+b64
`),
				},
			},
			expectedFiles: map[string]string{
				"file1.txt":     "some text in here",
				"dir/file2.txt": "some other text in here",
			},
		},
		{
			name: "Extract and save multi-line content using non mime cloud-init",
			attachments: []MimeAttachment{
				{
					ContentType: "text/cloud-config",
					Content: []byte(`#cloud-config
"hostname": "some hostname is here"
"runcmd":
- "update-ca-trust"
"write_files":
- "content": |
     ---
     document: some-defaults
  "path": "/etc/dnf/modules.defaults.d/postgresql.yaml"
`),
				},
			},
			expectedFiles: map[string]string{
				"/etc/dnf/modules.defaults.d/postgresql.yaml": "---\ndocument: some-defaults\n",
			},
		},
		{
			name: "Extract and save multi-line content",
			attachments: []MimeAttachment{
				{
					ContentType: "text/cloud-config",
					Content: []byte(`---
write_files:
  - path: multiline.txt
    content: |
      This is line 1.
      This is line 2.
      This is line 3.
`),
				},
			},
			expectedFiles: map[string]string{
				"multiline.txt": "This is line 1.\nThis is line 2.\nThis is line 3.\n",
			},
		},
		{
			name:          "No attachments",
			attachments:   []MimeAttachment{},
			expectedFiles: map[string]string{},
		},
		{
			name: "ReadCloudConfigFrom returns an error",
			attachments: []MimeAttachment{
				{
					ContentType: "text/cloud-config",
					Content:     []byte("invalid content"),
				},
			},
			expectedError: fmt.Errorf("failed to extract cloud config write files:"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("could not create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)
			for _, attachment := range tt.attachments {
				err = ExtractCloudConfig(attachment, tempDir)

				if tt.expectedError != nil {
					assert.True(t, strings.Contains(err.Error(), tt.expectedError.Error()), "error should contain: %s", tt.expectedError.Error())
					return
				}
				assert.NoError(t, err)
			}

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
