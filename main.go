package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	outputDir, input, err := parseFlags()
	if err != nil {
		return err
	}

	attachments, err := ExtractMimeAttachmentsFromBytes([]byte(input))
	if err != nil {
		return fmt.Errorf("failed to extract mime attachments: %w", err)
	}

	err = ExtractCloudConfig(attachments, outputDir)
	if err != nil {
		return fmt.Errorf("failed to save write files: %w", err)
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
