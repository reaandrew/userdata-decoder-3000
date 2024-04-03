package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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
		decoded, err := decode(input.Data)
		if err != nil {
			fmt.Println(fmt.Sprintf("There was an error decoding the base64 content %v", err))
			outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
			fullPath := filepath.Join(outputPath, "userdata")
			err := os.MkdirAll(path.Dir(fullPath), 0755)
			if err != nil {
				return fmt.Errorf("error creating output directories: %w", err)
			}
			err = os.WriteFile(fullPath, input.Data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}
		} else {
			attachments, err := ExtractMimeAttachmentsFromBytes(decoded)
			if err != nil {
				outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
				fullPath := filepath.Join(outputPath, "userdata")
				err := os.MkdirAll(path.Dir(fullPath), 0755)
				if err != nil {
					return fmt.Errorf("error creating output directories: %w", err)
				}
				err = os.WriteFile(fullPath, decoded, 0644)
				if err != nil {
					return fmt.Errorf("error writing file: %w", err)
				}
			} else {
				outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)

				err = ExtractCloudConfig(attachments, outputPath)
				if err != nil {
					return fmt.Errorf("failed to save write files: %w", err)
				}
			}
		}
	}

	return nil
}
