package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	diffStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	borderStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true)
	helpKeyColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("35")).Bold(true)
	helpDescColor   = lipgloss.NewStyle().Faint(true)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	normalLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
)

func Render(m Model) string {

	if len(m.urls) == 0 && !m.loading {
		return borderStyle.Render("No URLs found. Press 'q' to quit.")
	}

	var leftPanel strings.Builder
	leftPanel.WriteString(titleStyle.Render("Tracked URLs\n\n"))
	for i, url := range m.urls {
		if i == m.selectedIdx {
			leftPanel.WriteString(selectedStyle.Render("> "+url) + "\n")
		} else {
			leftPanel.WriteString(normalLineStyle.Render("  "+url) + "\n")
		}
	}
	leftView := borderStyle.Width(m.width/2 - 2).Render(leftPanel.String())

	var rightPanel strings.Builder
	if m.showDiff {
		rightPanel.WriteString(titleStyle.Render("Diff View\n\n"))
		rightPanel.WriteString(diffStyle.Render(m.diffContent) + "\n")
	} else {
		if m.loading {
			rightPanel.WriteString(titleStyle.Render("Processing...\n\n"))
			rightPanel.WriteString("Please wait...")
		} else {
			rightPanel.WriteString(titleStyle.Render("Instructions\n\n"))
			rightPanel.WriteString(renderInstructions())
		}
	}
	rightView := borderStyle.Width(m.width/2 - 2).Render(rightPanel.String())

	row := lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)

	statusBar := ""
	if m.statusMsg != "" {
		statusBar = "\n" + m.statusMsg
	}

	return row + statusBar
}

func renderInstructions() string {
	return fmt.Sprintf("%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
		helpKeyColor.Render("↑/↓ or j/k:"), helpDescColor.Render("Select URL"),
		helpKeyColor.Render("d:"), helpDescColor.Render("Show diff for selected URL"),
		helpKeyColor.Render("r:"), helpDescColor.Render("Re-check selected URL"),
		helpKeyColor.Render("q:"), helpDescColor.Render("Quit TUI"),
		helpKeyColor.Render("esc:"), helpDescColor.Render("Quit TUI"),
	)
}
