package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Terminal has some useful information about our terminal screen.
// We'll move it to its own file/pkg later.
type Terminal struct {
	Width  int
	Height int
}

// Model is the state of our app.
type Model struct {
	// Options for our application.
	// TODO: We need to remove this options from the model, and set the appropriate
	// values in the `initialModel` function.
	Options
	// Data contains our list of notes.
	Data
	// Glyphs give us the runes we need according to the selected options.
	Glyphs
	// Terminal has useful info about our screen.
	Terminal
	// Quit means quit (but we'll comment any public symbols anyway haha).
	Quit bool
}

func initialModel(opts Options) Model {
	glyphs := GetGlyphs(opts.Graphics)
	return Model{Glyphs: glyphs, Options: opts}
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

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

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

		case 'm':
			// indent [m]ore (increase tab size).
			// max: 8 columns.
			if m.TabSize < 8 {
				m.TabSize++
			}

		case 'l':
			// indent [l]ess (decrease tab size).
			// min: 1 column.
			if m.TabSize > 1 {
				m.TabSize--
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

	content.WriteString(m.Styled())

	// TODO: There are bubbles for this.
	// TODO: We don't need to show the keybindings in the UI (maybe just a hint for help).
	// Maybe show a `Usage` message at first launch.  We can add `Help` to the command with `--help`.
	content.WriteString("\n\n\x1b[2m j : down\n k : up\n t : toggle collapsed\n o : open all\n c : close all\n q : quit\n\x1b[0m")

	fmt.Fprintf(&content, "\n\x1b[32m%dx%d\n", m.Width, m.Height)

	view := tea.NewView(content.String())
	view.AltScreen = true

	return view
}

func (m Model) ToString() string {

	var s strings.Builder

	type list struct {
		notes []Note
		pos   int
	}

	// lists will be a stack of (sub)lists of notes.
	lists := []list{
		{
			notes: m.Notes,
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
			if m.Selected == subList.pos {
				cursor = m.Cursor
			}
		}

		// Mark elements with sub-lists with a '+'
		more := "  "
		if len(currentNote.Items) > 0 {
			more = m.More
		}

		tab := strings.Repeat(" ", m.TabSize)
		indent := strings.Repeat(tab, len(lists)-1)

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

// Test: using Lipgloss.
func (m Model) Styled() string {
	leftPane := lipgloss.NewStyle().
		Padding(1).
		Width(m.Width / 2)

	rightPane := lipgloss.NewStyle().
		Padding(0, 2).
		Width(m.Width - (m.Width / 2)).
		Foreground(lipgloss.Yellow).
		Border(lipgloss.RoundedBorder())

	highlight := lipgloss.NewStyle().
		Foreground(lipgloss.Blue).
		Reverse(true)

	var leftContent strings.Builder
	var rightContent strings.Builder

	for i, note := range m.Notes {
		more := "  "
		if len(note.Items) > 0 {
			more = m.More
		}

		line := fmt.Sprintf("%s %s", more, note.Summary)

		if i == m.Selected {
			line = highlight.Render(line)
			if len(note.Items) > 0 {
				for _, item := range note.Items {
					fmt.Fprintf(&rightContent, "- %s\n", item.Summary)
				}
			}
		}

		leftContent.WriteString(line + "\n")
	}

	if rightContent.Len() == 0 {
		rightContent.WriteString("empty")
	}

	left := leftPane.Render(leftContent.String())
	right := rightPane.Render(strings.TrimSpace(rightContent.String()))

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
