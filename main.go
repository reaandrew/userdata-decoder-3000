package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	output, input, err := parseFlags()
	if err != nil {
		return err
	}

	attachments, err := ExtractMimeAttachmentsFromBytes([]byte(input))
	if err != nil {
		return fmt.Errorf("failed to extract mime attachments: %w", err)
	}

	for _, attachment := range attachments {
		if strings.Contains(attachment.ContentType, "text/cloud-config") {
			cloudConfig, err := ReadCloudConfigFrom(attachment, output)
			if err != nil {
				return fmt.Errorf("failed to extract cloud config write files: %w", err)
			}

			if err := saveWriteFiles(cloudConfig.WriteFiles); err != nil {
				return fmt.Errorf("failed to save write files: %w", err)
			}
		}
	}

	return nil
}

func parseFlags() (output string, input string, err error) {
	flag.StringVar(&output, "o", "output", "Specify the output directory within your working directory.")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		return "", "", errors.New("input file is required")
	}

	return output, args[0], nil
}

func saveWriteFiles(writeFiles []WriteFile) error {
	for _, file := range writeFiles {
		err := os.MkdirAll(path.Dir(file.Path), 0755)
		if err != nil {
			return fmt.Errorf("error creating output directories: %w", err)
		}
		err = os.WriteFile(file.Path, []byte(file.Content), 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
	}
	return nil
}
