package main

type DataOutputPair struct {
	Data      []byte
	OutputDir string
}

type DataProvider interface {
	FetchData() ([]DataOutputPair, error)
}
