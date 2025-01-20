package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/copyleftdev/snaptrack/pkg/store"
)

func StartProgram(db store.DBInterface) error {
	m := NewModel(db)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Println("Error running TUI:", err)
		return err
	}
	return nil
}
