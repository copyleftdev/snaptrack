package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/copyleftdev/snaptrack/pkg/capture"
	"github.com/copyleftdev/snaptrack/pkg/diff"
	"github.com/copyleftdev/snaptrack/pkg/snapshot"
	"github.com/copyleftdev/snaptrack/pkg/store"
)

type Model struct {
	db store.DBInterface

	urls        []string
	selectedIdx int
	diffContent string
	showDiff    bool
	statusMsg   string
	loading     bool
	width       int
	height      int
}

func NewModel(db store.DBInterface) Model {
	return Model{
		db:   db,
		urls: []string{},
	}
}

func (m Model) Init() tea.Cmd {

	return loadURLsCmd(m.db)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case loadURLsMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Error loading URLs: %v", msg.err)
			return m, nil
		}
		m.urls = msg.urls
		return m, nil

	case diffMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Diff error: %v", msg.err)
			m.showDiff = false
			return m, nil
		}
		m.diffContent = msg.diffText
		m.showDiff = true
		m.loading = false
		return m, nil

	case checkMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Check error: %v", msg.err)
		} else {
			m.statusMsg = msg.info
		}
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			m.showDiff = false
			m.diffContent = ""
			return m, nil

		case "down", "j":
			if m.selectedIdx < len(m.urls)-1 {
				m.selectedIdx++
			}
			m.showDiff = false
			m.diffContent = ""
			return m, nil

		case "d":
			if len(m.urls) > 0 {
				selectedURL := m.urls[m.selectedIdx]
				m.showDiff = true
				m.diffContent = "Loading diff..."
				return m, generateDiffCmd(m.db, selectedURL)
			}

		case "r":
			if len(m.urls) > 0 {
				selectedURL := m.urls[m.selectedIdx]
				m.loading = true
				m.statusMsg = ""
				return m, recheckURLCmd(m.db, selectedURL)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return Render(m)
}

type loadURLsMsg struct {
	urls []string
	err  error
}

func loadURLsCmd(db store.DBInterface) tea.Cmd {
	return func() tea.Msg {
		urls, err := getDistinctURLs(db)
		return loadURLsMsg{urls, err}
	}
}

type diffMsg struct {
	diffText string
	err      error
}

func generateDiffCmd(db store.DBInterface, url string) tea.Cmd {
	return func() tea.Msg {
		snaps, err := getSnapshotsForURL(db, url)
		if err != nil {
			return diffMsg{"", err}
		}
		if len(snaps) < 2 {
			return diffMsg{"No diff available (need at least 2 snapshots).", nil}
		}
		oldSnap := snaps[1].HTML
		newSnap := snaps[0].HTML
		diffText := diff.GenerateDiff(oldSnap, newSnap, true)
		return diffMsg{diffText, nil}
	}
}

type checkMsg struct {
	info string
	err  error
}

func recheckURLCmd(db store.DBInterface, url string) tea.Cmd {
	return func() tea.Msg {

		html, reqHeaders, respHeaders, statusCode, err := capture.CaptureHTML(url, 15*time.Second)
		if err != nil {
			return checkMsg{"", err}
		}

		err = snapshot.StoreOrUpdateSnapshot(db, url, html, statusCode, reqHeaders, respHeaders)
		if err != nil {
			return checkMsg{"", err}
		}
		return checkMsg{fmt.Sprintf("Checked %s successfully.", url), nil}
	}
}

func getDistinctURLs(db store.DBInterface) ([]string, error) {

	type distinctURLoader interface {
		GetDistinctURLs() ([]string, error)
	}
	if loader, ok := db.(distinctURLoader); ok {
		return loader.GetDistinctURLs()
	}
	return nil, fmt.Errorf("db does not implement GetDistinctURLs")
}

func getSnapshotsForURL(db store.DBInterface, url string) ([]snapshot.Snapshot, error) {
	type snapshotsForURL interface {
		GetSnapshotsForURL(string) ([]snapshot.Snapshot, error)
	}
	if snapper, ok := db.(snapshotsForURL); ok {
		return snapper.GetSnapshotsForURL(url)
	}
	return nil, fmt.Errorf("db does not implement GetSnapshotsForURL")
}
