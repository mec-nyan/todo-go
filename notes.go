package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

// Collapse will hide the list of items.
// NOTE: Should we return another note instead of modifying it?
func (n *Note) Collapse() {
	n.Show = false
}

func (n *Note) Expand() {
	n.Show = true
}

func (n *Note) CollapseAll() {
	n.Collapse()
	// Again, a dirty recursive solution.  It's acceptable since we won't have that
	// many levels of recursion.
	// NOTE: There's no need to check the length of n.Items.
	// If it is empty, this loop will be skipped.
	for _, item := range n.Items {
		item.CollapseAll()
	}
}

func (n *Note) ExpandAll() {
	n.Expand()
	// Again, a dirty recursive solution.  It's acceptable since we won't have that
	// many levels of recursion.
	// NOTE: There's no need to check the length of n.Items.
	// If it is empty, this loop will be skipped.
	for _, item := range n.Items {
		item.Expand()
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

func (d *Data) Next() {
	d.Selected++
	if d.Selected == len(d.Notes) {
		d.Selected = 0
	}
}

func (d *Data) Previous() {
	if d.Selected > 0 {
		d.Selected--
	} else {
		d.Selected = len(d.Notes) - 1
	}
}

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

func (d *Data) Collapse() {
	setCollapsed(d.Notes, true)
}

func (d *Data) Expand() {
	setCollapsed(d.Notes, false)
}

func (d *Data) ToString() string {

	var s strings.Builder

	type list struct {
		notes []Note
		pos   int
	}

	// lists will be a stack of (sub)lists of notes.
	lists := []list{
		{
			notes: d.Notes,
			pos:   0,
		},
	}

	for {
		// Get the last list on the stack.
		subList := &lists[len(lists)-1]

		if subList.pos == len(subList.notes) {
			if len(lists) == 1 {
				// We're done here.
				break
			}
			// Go back one level.
			lists = lists[:len(lists)-1]

			// Re-check the condition.
			continue
		}

		// Ok, we've got an element to add.
		currentNote := subList.notes[subList.pos]

		// For now, navigate only in the main list.
		// TODO: We also need to navegate the sub-lists.
		cursor := "  "
		if len(lists) == 1 {
			// We're on the main list.
			if d.Selected == subList.pos {
				cursor = ">>"
			}
		}

		// Mark elements with sub-lists with a '+'
		more := "  "
		if len(currentNote.Items) > 0 {
			more = "++"
		}

		indent := strings.Repeat("    ", len(lists)-1)

		fmt.Fprintf(&s, " %s %s%s %s\n", cursor, indent, more, currentNote.Summary)

		// Mark this element as done by moving to the next pos.
		subList.pos++

		// If this element has a sub-list, and it's not collapsed, push it to the stack.
		if len(currentNote.Items) > 0 && currentNote.Show {
			lists = append(lists, list{notes: currentNote.Items, pos: 0})
		}
	}

	return s.String()
}
