package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Note represent a single Note.  It can have a 'Summary' (the Note itself) and optionally,
// a list of items.  Each item is itself another Note, so we can have nested lists.
type Note struct {
	// Summary is the main content of the note (i.e. 'Buy groceries').
	Summary string `json:"summary"`
	// Items contains an optional list of elements (i.e. 'lettuce, cucumber, tomatoes, ...').
	// Each item, in turn, can have a list of sub-elements, etc.
	Items []Note `json:"items"`
	// TODO: Maybe add more fields (i.e. 'title', 'description', 'summary', etc).
}

// Data is the content of our `Data.json` file..
type Data struct {
	// Notes is an array of `Note`s.
	Notes []Note `json:"notes"`
}

// LoadNotes reads a `json` file and loads its content.
// It should be a valid list of notes.
func LoadNotes(filename string) (*Data, error) {
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
