package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Model is the state of our app.
type Model struct {
	// Options for our application.
	Options
	// Data contains our list of notes.
	Data
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
		data, err := LoadNotes(m.SaveFile)
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
			m.Next()

		case 'k':
			m.Previous()

		case 'c':
			// [c]lose all.
			m.Collapse()

		case 'o':
			// [o]pen all.
			m.Expand()

		case 'e':
			// [e]xpand this item, one level.
			// '+' will show items with more levels/items in it.
			note := m.GetSelected()
			if note.Show {
				note.Collapse()
			} else {
				note.Expand()
			}

		case 'r':
			// [r]ecursively expand this item (all levels)
			note := m.GetSelected()
			if note.Show {
				note.CollapseAll()
			} else {
				note.ExpandAll()
			}

		case 'q', tea.KeyEscape:
			m.Quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// TODO: The view is very primitive now.  We'll improve it later!
// View fulfills tea.Model.
func (m Model) View() tea.View {

	if m.Quit {
		return tea.NewView("Bye!")
	}

	var content strings.Builder

	content.WriteString(" Notes:\n\n")

	content.WriteString(m.ToString())

	// TODO: There are bubbles for this.
	// TODO: We don't need to show the keybindings in the UI (maybe just a hint for help).
	// Maybe show a `Usage` message at first launch.  We can add `Help` to the command with `--help`.
	content.WriteString("\n\n\x1b[2m j : down\n k : up\n t : toggle collapsed\n o : open all\n c : close all\n q : quit\n")

	view := tea.NewView(content.String())
	view.AltScreen = true

	return view
}
