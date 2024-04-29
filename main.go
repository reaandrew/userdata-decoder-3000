package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

var Log = logrus.New()
var Version string

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
	flag.BoolVar(&config.version, "version", false, "Output the current version")
	flag.BoolVar(&config.verbose, "v", false, "Output debug information from the process")
	flag.BoolVar(&config.verbose, "verbose", false, "Output debug information from the process")
	flag.StringVar(&config.outputDir, "o", "output", "Specify the output directory within your working directory.")
	flag.StringVar(&config.outputDir, "output-dir", "output", "Specify the output directory within your working directory.")

	flag.Usage = func() {
		fmt.Println("Usages:")
		fmt.Println("  udd [OPTIONS]             : Specify options for data provider and output directory.")
		fmt.Println("  udd [CONTENT_TO_DECODE]   : Specify the content to decode as the first argument.")
		fmt.Println("\nOptions:")
		fmt.Println("  -o, --output-dir: Specify the output directory within your working directory. (default \"output\")")
		fmt.Println("  -p, --provider:   Specify the data provider (e.g., aws).")
		fmt.Println("  -v, --verbose:   Output debug information from the process")
		fmt.Println("  	   --version:   OOutput the current version")
		os.Exit(0)
	}

	flag.Parse()

	args := flag.Args()

	Log.Formatter = new(logrus.JSONFormatter)
	Log.Out = os.Stdout

	if config.version {
		fmt.Println("Version: " + Version)
		os.Exit(0)
	}

	if config.verbose {
		Log.Level = logrus.DebugLevel
	} else {
		Log.Level = logrus.ErrorLevel
	}

	if config.providerKey == "" && len(args) < 1 {
		return Config{}, errors.New("Either supply content directly or use a provider")
	}
	config.args = args

	return config, nil
}
