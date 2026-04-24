package main

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
)

type model int

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.Code {

		case 'j', '+':
			m++

		case 'k', '-':
			m--

		case 'q', tea.KeyEscape:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	out := "Counter App\n\n"
	out += fmt.Sprintf("    Count: %d\n\n", m)
	out += "\x1b[2mj/+: increment; k/-: decrement; q/esc: quit\x1b[0m\n"
	return tea.NewView(out)
}

func main() {
	var m model = 0

	app := tea.NewProgram(m)
	if _, err := app.Run(); err != nil {
		log.Fatalf("Ooooooooooops! (%v)", err)
	}
}
