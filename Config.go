package main

import (
	"fmt"
)

var providers = map[string]DataProvider{
	"aws": NewDefaultAwsProvider(),
}

type Config struct {
	outputDir   string
	providerKey string
	verbose     bool
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
			return nil, fmt.Errorf("unknown provider: %s", config.providerKey)
		}
		return provider, nil
	}
}
