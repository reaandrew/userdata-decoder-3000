package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type InputProcessor struct {
	config Config
}

func (inputProcessor InputProcessor) writePlainUserDataFile(input DataOutputPair) error {
	outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
	fullPath := filepath.Join(outputPath, "userdata")
	err := os.MkdirAll(path.Dir(fullPath), 0755)

	if err != nil {
		return fmt.Errorf("error creating output directories: %w", err)
	}
	err = os.WriteFile(fullPath, input.Data, 0755)
	if err != nil {
		return fmt.Errorf("errir writing file: %w", err)
	}
	return nil
}

func (inputProcessor InputProcessor) Process(inputs []DataOutputPair) error {
	for _, input := range inputs {
		outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
		fullPath := filepath.Join(outputPath, "raw")
		err := os.MkdirAll(path.Dir(outputPath), 0755)
		if err != nil {
			return fmt.Errorf("error creating output directories: %w", err)
		}
		err = os.WriteFile(fullPath, input.Data, 0644)

		decoded, err := decode(input.Data)
		rawPath := filepath.Join(outputPath, "raw_base64Decoded")
		rawErr := os.MkdirAll(path.Dir(outputPath), 0755)
		if rawErr != nil {
			return fmt.Errorf("error creating output directories: %w", rawErr)
		}
		rawErr = os.WriteFile(rawPath, decoded, 0644)

		if err != nil {
			fmt.Println(fmt.Sprintf("There was an error decoding the base64 content %v", err))
			fullPath := filepath.Join(outputPath, "userdata")
			err := os.MkdirAll(path.Dir(outputPath), 0755)
			if err != nil {
				return fmt.Errorf("error creating output directories: %w", err)
			}
			err = os.WriteFile(fullPath, input.Data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file after failing to decode base64: %w", err)
			}
		} else {
			attachments, err := ExtractMimeAttachmentsFromBytes(decoded)
			if err != nil {
				err := os.MkdirAll(path.Dir(outputPath), 0755)
				fullPath := filepath.Join(outputPath, "userdata")
				if err != nil {
					return fmt.Errorf("error creating output directories: %w", err)
				}
				err = os.WriteFile(fullPath, decoded, 0644)
				if err != nil {
					return fmt.Errorf("error writing file after failing to extract mime attachments: %w", err)
				}
			} else {
				for _, attachment := range attachments {
					if strings.Contains(attachment.ContentType, "text/cloud-config") ||
						strings.Contains(string(attachment.Content), "#cloud-config") {
						err = ExtractCloudConfig(attachment, outputPath)
						if err != nil {
							return fmt.Errorf("failed to save write files: %w", err)
						}
					} else {
						fullPath := filepath.Join(outputPath, attachment.Filename)
						err := os.MkdirAll(path.Dir(outputPath), 0755)
						if err != nil {
							return fmt.Errorf("error creating output directories: %w", err)
						}
						decodedContent, err := decode(attachment.Content)
						err = os.WriteFile(fullPath, decodedContent, 0644)
					}
				}
			}
		}
	}

	return nil
}
