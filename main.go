package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

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
		log.Fatalf("failed to the provider")
	}

	inputs, err := provider.FetchData()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	return InputProcessor{
		config: config,
	}.Process(inputs)
}

func parseFlags() (config Config, err error) {
	flag.StringVar(&config.providerKey, "p", "", "Specify the data provider (e.g., aws).")
	flag.StringVar(&config.providerKey, "provider", "", "Specify the data provider (e.g., aws).")
	flag.StringVar(&config.outputDir, "o", "output", "Specify the output directory within your working directory.")
	flag.StringVar(&config.outputDir, "output-dir", "output", "Specify the output directory within your working directory.")

	flag.Usage = func() {
		fmt.Println("Usages:")
		fmt.Println("  cloud-init-decoder [OPTIONS]             : Specify options for data provider and output directory.")
		fmt.Println("  cloud-init-decoder [CONTENT_TO_DECODE]   : Specify the content to decode as the first argument.")
		fmt.Println("\nOptions:")
		fmt.Println("  -o, --output-dir: Specify the output directory within your working directory. (default \"output\")")
		fmt.Println("  -p, --provider:   Specify the data provider (e.g., aws).")
		os.Exit(0)
	}

	flag.Parse()

	args := flag.Args()

	if config.providerKey == "" && len(args) < 1 {
		return Config{}, errors.New("Either supply content directly or use a provider")
	}
	config.args = args

	return config, nil
}
