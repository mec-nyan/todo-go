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
	// Show determines if the list of items is hidden.
	// NOTE: Should we save this state to file?  If so, the app will preserve its UI's state between
	// launches.  Maybe it's a good idea, maybe not.  We can add a configuration option for this.
	Show bool
	// TODO: Maybe add more fields (i.e. 'title', 'description', 'summary', etc).
}

// Collapse the note to hide its items.
// NOTE: Should we return another note instead of modifying it?
func (n *Note) Collapse() {
	n.Show = false
}

// Expand the note to show its items.
func (n *Note) Expand() {
	n.Show = true
}

// CollapseAll hides the note's items recursively.
func (n *Note) CollapseAll() {
	n.Collapse()
	// Again, a dirty recursive solution.  It's acceptable since we won't have that
	// many levels of recursion.
	// NOTE: There's no need to check the length of n.Items.
	// If it is empty, this loop will be skipped.
	for i := range n.Items {
		n.Items[i].CollapseAll()
	}
}

// ExpandAll shows the note's items recursively.
func (n *Note) ExpandAll() {
	n.Expand()
	// Again, a dirty recursive solution.  It's acceptable since we won't have that
	// many levels of recursion.
	// NOTE: There's no need to check the length of n.Items.
	// If it is empty, this loop will be skipped.
	for i := range n.Items {
		n.Items[i].ExpandAll()
	}
}

// Data is the content of our `Data.json` file..
type Data struct {
	// Notes is an array of `Note`s.
	Notes []Note `json:"notes"`
	// Selected is the index of the currently selected note.
	Selected int
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

// Select the next item in the list (with wraparound).
func (d *Data) Next() {
	d.Selected++
	if d.Selected == len(d.Notes) {
		d.Selected = 0
	}
}

// Select the previous item in the list (with wraparound).
func (d *Data) Previous() {
	if d.Selected > 0 {
		d.Selected--
	} else {
		d.Selected = len(d.Notes) - 1
	}
}

// Get a reference to the current item.
func (d *Data) GetSelected() *Note {
	return &d.Notes[d.Selected]
}

func setCollapsed(notes []Note, collapsed bool) {
	// For now, an easy solution using recursion.  We won't have that many levels anyway.
	// TODO: Write this without recursion.
	for i := range notes {
		note := &notes[i]
		note.Show = !collapsed
		if len(note.Items) > 0 {
			setCollapsed(note.Items, collapsed)
		}
	}
}

// Collapse every item of the list (recursively).
func (d *Data) Collapse() {
	setCollapsed(d.Notes, true)
}

// Expand every item of the list (recursively).
func (d *Data) Expand() {
	setCollapsed(d.Notes, false)
}
