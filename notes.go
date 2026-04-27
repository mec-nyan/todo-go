package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Note struct {
	Summary string `json:"summary"`
	Items   []Note `json:"items"`
}

type Data struct {
	Notes []Note `json:"notes"`
}

func loadNotes(filename string) (*Data, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load file %s: %w", filename, err)
	}

	var data Data
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &data, nil
}
