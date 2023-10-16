package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var providers = map[string]DataProvider{
	"aws": AWSProvider{},
}

type Config struct {
	outputDir   string
	providerKey string
	args        []string
}

func (config Config) getProvider() (DataProvider, error) {
	if config.providerKey == "" {
		input := config.args[0]
		return CommandLineProvider{Input: input}, nil
	} else {
		var ok bool
		provider, ok := providers[config.providerKey]
		if !ok {
			fmt.Println("Unknown provider:", config.providerKey)
			os.Exit(1)
		}
		return provider, nil
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := parseFlags()

	if err != nil {
		return err
	}
	provider, err := config.getProvider()
	if err != nil {
		return err
	}

	inputs, err := provider.FetchData()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	for _, input := range inputs {
		attachments, err := ExtractMimeAttachmentsFromBytes(input.Data)
		if err != nil {
			return fmt.Errorf("failed to extract mime attachments: %w", err)
		}

		outputPath := filepath.Join(config.outputDir, input.OutputDir)

		err = ExtractCloudConfig(attachments, outputPath)
		if err != nil {
			return fmt.Errorf("failed to save write files: %w", err)
		}
	}

	return nil
}

func parseFlags() (config Config, err error) {
	flag.StringVar(&config.providerKey, "provider", "", "Specify the data provider (e.g., aws).")
	flag.StringVar(&config.outputDir, "o", "output", "Specify the output directory within your working directory.")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		return Config{}, errors.New("Either supply content directly or use a provider")
	}
	config.args = args

	return config, nil
}
