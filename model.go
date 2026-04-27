package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Model is the state of our app.
type Model struct {
	// Data contains our list of notes.
	Data
	// Current indicates the currently selected note.
	Current int
	// Collapse indicates wether to show all the sub-items or just the main list.
	Collapse bool
	// Quit means quit (but we'll comment any public symbols anyway haha).
	Quit bool
}

// FileLoader is the "message" that contains the notes saved to file,
// or an error if we couldn't load it.
type FileLoader struct {
	Data  *Data
	Error error
}

// Init will load the notes from file, if any.
func (m Model) Init() tea.Cmd {
	// TODO: File path should be set by configuration and/or command line argument.
	// TODO: If the file doesn't exist (i.e. first launch) we need to create it in the right place.
	return func() tea.Msg {
		data, err := LoadNotes("data.json")
		return FileLoader{
			Data:  data,
			Error: err,
		}
	}

}

// Update fulfills tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case FileLoader:
		// TODO: Handle errors properly.
		if msg.Error != nil {
			return m, tea.Quit
		}

		m.Data = *msg.Data
		return m, nil

	case tea.KeyPressMsg:

		switch msg.Code {

		case 'j':
			m.Current++
			if m.Current == len(m.Notes) {
				m.Current = 0
			}

		case 'k':
			m.Current--
			if m.Current < 0 {
				m.Current = len(m.Notes) - 1
			}

		case 'c':
			m.Collapse = true

		case 'o':
			m.Collapse = false

		case 'q', tea.KeyEscape:
			m.Quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View fulfills tea.Model.
func (m Model) View() tea.View {

	if m.Quit {
		return tea.NewView("Bye!")
	}

	var view strings.Builder

	// view.WriteString("Notes:\n\n")

	// for i, note := range m.Notes {
	// 	cursor := " "
	// 	if i == m.Current {
	// 		cursor = ">"
	// 	}

	// 	more := " "
	// 	if len(note.Items) > 0 {
	// 		more = "+"
	// 	}

	// 	fmt.Fprintf(&view, "  %s %s %s\n", cursor, more, note.Summary)
	// }

	view.WriteString(showNotes(m.Notes, 0))

	// TODO: There are bubbles for this.
	view.WriteString("\n\n\x1b[2m  j : down - k : up - q : quit\n")

	return tea.NewView(view.String())
}

// showNotes 'formats' our notes in the way we want (i.e. show or collapse items).
func showNotes(notes []Note, _ int) string {
	var s strings.Builder

	type list struct {
		notes []Note
		pos   int
	}

	// lists will be a stack of (sub)lists of notes.
	lists := []list{
		{
			notes: notes,
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

		indent := strings.Repeat("    ", len(lists)-1)

		fmt.Fprintf(&s, "%s- %s.\n", indent, currentNote.Summary)

		// Mark this element as done by moving to the next pos.
		subList.pos++

		// If this element has a sub-list, push it to the stack.
		if len(currentNote.Items) > 0 {
			lists = append(lists, list{notes: currentNote.Items, pos: 0})
		}
	}

	return s.String()
}
