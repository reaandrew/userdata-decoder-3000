package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

type WriteFile struct {
	Path     string `yaml:"path"`
	Encoding string `yaml:"encoding"`
	Content  string `yaml:"content"`
}

type CloudConfig struct {
	WriteFiles []WriteFile `yaml:"write_files"`
}

func ReadCloudConfigFrom(attachment MimeAttachment, baseDir string) (CloudConfig, error) {
	if !strings.Contains(attachment.ContentType, "text/cloud-config") {
		return CloudConfig{}, fmt.Errorf("not a cloud-config content type")
	}

	var config CloudConfig

	err := yaml.Unmarshal(attachment.Content, &config)
	if err != nil {
		return CloudConfig{}, err
	}

	var writeFiles []WriteFile
	for _, file := range config.WriteFiles {
		fullPath := filepath.Join(baseDir, file.Path)
		content, err := decode([]byte(file.Content))
		if err != nil {
			return CloudConfig{}, err
		}

		writeFiles = append(writeFiles, WriteFile{
			Path:     fullPath,
			Encoding: file.Encoding,
			Content:  string(content),
		})
	}
	return CloudConfig{WriteFiles: writeFiles}, err
}
