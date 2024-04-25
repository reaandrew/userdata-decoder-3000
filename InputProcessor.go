package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InputProcessor struct {
	config Config
}

func (inputProcessor InputProcessor) writeFile(input DataOutputPair, filename string) error {
	outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
	writeFilePath := filepath.Join(outputPath, filename)

	return inputProcessor.writeFileWithData(input.Data, writeFilePath)
}

func (inputProcessor InputProcessor) writeFileWithData(data []byte, filename string) error {
	pathToMake := filepath.Dir(filename)
	err := os.MkdirAll(pathToMake, 0755)

	if err != nil {
		return fmt.Errorf("error creating output directories: %w", err)
	}
	err = os.WriteFile(filename, data, 0755)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func (inputProcessor InputProcessor) Process(inputs []DataOutputPair) error {
	Log.Debugln("Processing inputs")
	for _, input := range inputs {
		outputPath := filepath.Join(inputProcessor.config.outputDir, input.OutputDir)
		err := inputProcessor.writeFile(input, "raw")
		if err != nil {
			return err
		}

		Log.WithField("decoded_data", string(input.Data)).Debug("Encoded Data")
		decoded, err := decode(input.Data)
		rawErr := inputProcessor.writeFileWithData(decoded, filepath.Join(outputPath, "raw_base64Decoded"))
		if rawErr != nil {
			return rawErr
		}
		Log.WithField("decoded_data", string(decoded)).Debug("Decoded Data")
		if err != nil {
			fmt.Println(fmt.Sprintf("There was an error decoding the base64 content %v", err))
			err := inputProcessor.writeFile(input, "userdata")
			if err != nil {
				return fmt.Errorf("error writing file after failing to decode base64: %w", err)
			}
		} else {
			attachments, err := ExtractMimeAttachmentsFromBytes(decoded)
			if err != nil {
				Log.Errorln("Error Extracting Mime Attachments From Bytes", string(decoded))

				err := inputProcessor.writeFile(input, "userdata")
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
						decodedContent, err := decode(attachment.Content)
						err = inputProcessor.writeFileWithData(decodedContent, fullPath)
						if err != nil {
							return err
						}

					}
				}
			}
		}
	}

	return nil
}
