package main

type CommandLineProvider struct {
	Input string
}

func (p CommandLineProvider) FetchData() ([]DataOutputPair, error) {
	return []DataOutputPair{
		{
			Data:      []byte(p.Input),
			OutputDir: "",
		},
	}, nil
}
