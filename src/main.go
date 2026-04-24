package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type Note struct {
	Summary string `json:"summary"`
	Items   []Note `json:"items"`
}

type Data struct {
	Notes []Note `josn:"notes"`
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

type model struct {
	Data
	Current  int
	Collapse bool
	Quit     bool
}

type fileLoader struct {
	Data  *Data
	Error error
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		// TODO: Use configuration or command line arguments for the file path.
		data, err := loadNotes("data.json")
		return fileLoader{
			Data:  data,
			Error: err,
		}
	}

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case fileLoader:
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

func (m model) View() tea.View {

	if m.Quit {
		return tea.NewView("Bye!")
	}

	var view strings.Builder

	view.WriteString("Notes:\n\n")

	for i, note := range m.Notes {
		cursor := " "
		if i == m.Current {
			cursor = ">"
		}

		more := " "
		if len(note.Items) > 0 {
			more = "+"
		}

		fmt.Fprintf(&view, "  %s %s %s\n", cursor, more, note.Summary)
	}

	// TODO: There are bubbles for this.
	view.WriteString("\n\n\x1b[2m  j : down - k : up - q : quit\n")

	return tea.NewView(view.String())
}

func main() {

	app := tea.NewProgram(model{Collapse: true})

	if _, err := app.Run(); err != nil {
		log.Fatalf("Ooooooooooops! (%v)", err)
	}
}
