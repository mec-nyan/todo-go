package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {

	app := tea.NewProgram(Model{Collapse: true})

	if _, err := app.Run(); err != nil {
		log.Fatalf("Ooooooooooops! (%v)", err)
	}
}
