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
		outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)

		// Simplify directory creation and error handling
		if err := makeDirectory(outputPath); err != nil {
			return err
		}

		rawFullPath := filepath.Join(outputPath, "raw")
		if err := os.WriteFile(rawFullPath, input.Data, 0644); err != nil {
			return fmt.Errorf("error writing raw file: %w", err)
		}

		decoded, err := decode(input.Data)
		if rawErr := writeDataToFile(outputPath, "raw_base64decoded", decoded); err != nil {
			return rawErr
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("There was an error decoding the base64 content %v", err))
			// Since decoding failed, write the raw data to userdata file
			if err := writeDataToFile(outputPath, "userdata", input.Data); err != nil {
				return err
			}
		} else {
			attachments, err := ExtractMimeAttachmentsFromBytes(decoded)
			if err != nil {
				// Since MIME extraction failed, write the decoded data to userdata file
				if err := writeDataToFile(outputPath, "userdata", decoded); err != nil {
					return err
				}
			} else {
				if err := ExtractCloudConfig(attachments, outputPath); err != nil {
					return fmt.Errorf("failed to save write files: %w", err)
				}
				for _, attachment := range attachments {
					if err := writeDataToFile(outputPath, attachment.Filename, attachment.Content); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// makeDirectory creates the necessary directories
func makeDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("error creating output directories: %w", err)
	}
	return nil
}

// writeDataToFile creates the directory (if necessary), and writes the data to the specified filename
func writeDataToFile(outputPath, filename string, data []byte) error {
	fullPath := filepath.Join(outputPath, filename)
	if err := os.MkdirAll(path.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("error creating output directories: %w", err)
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
