package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	opts := parseArgs()

	app := tea.NewProgram(Model{Options: opts, Collapse: true})

	if _, err := app.Run(); err != nil {
		log.Fatalf("Ooooooooooops! (%v)", err)
	}
}
