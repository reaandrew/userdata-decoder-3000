package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
)

type WriteFile struct {
	Path     string `yaml:"path"`
	Encoding string `yaml:"encoding"`
	Content  string `yaml:"content"`
}

type CloudConfig struct {
	WriteFiles []WriteFile `yaml:"write_files"`
}

func (cloudConfig CloudConfig) SaveWriteFiles(outputDir string) error {
	for _, file := range cloudConfig.WriteFiles {
		fullPath := filepath.Join(outputDir, file.Path)
		err := os.MkdirAll(path.Dir(fullPath), 0755)
		if err != nil {
			return fmt.Errorf("error creating output directories: %w", err)
		}
		err = os.WriteFile(fullPath, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
	}
	return nil
}

func ReadCloudConfigFrom(attachment MimeAttachment) (CloudConfig, error) {
	var config CloudConfig

	err := yaml.Unmarshal(attachment.Content, &config)
	if err != nil {
		return CloudConfig{}, err
	}

	var writeFiles []WriteFile
	for _, file := range config.WriteFiles {
		content, err := decode([]byte(file.Content))
		if err != nil {
			return CloudConfig{}, err
		}

		writeFiles = append(writeFiles, WriteFile{
			Path:     file.Path,
			Encoding: file.Encoding,
			Content:  string(content),
		})
	}
	return CloudConfig{WriteFiles: writeFiles}, err

}

func ExtractCloudConfig(attachment MimeAttachment, outputDir string) error {
	cloudConfig, err := ReadCloudConfigFrom(attachment)
	if err != nil {
		return fmt.Errorf("failed to extract cloud config write files: %w", err)
	}

	if err := cloudConfig.SaveWriteFiles(outputDir); err != nil {
		return fmt.Errorf("failed to save write files: %w", err)
	}

	return nil
}
