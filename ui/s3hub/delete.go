package s3hub

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/ui"
)

type s3hubDeleteContentsModel struct {
	// quitting is true when the user has quit the application.
	quitting bool
}

func (m *s3hubDeleteContentsModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubDeleteContentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubDeleteContentsModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubDeleteContentsModel",
		ui.Subtle("j/k, up/down: select")+" | "+ui.Subtle("enter: choose")+" | "+ui.Subtle("q, esc: quit"))
}

type s3hubDeleteBucketModel struct {
	// quitting is true when the user has quit the application.
	quitting bool
}

func (m *s3hubDeleteBucketModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubDeleteBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubDeleteBucketModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubDeleteBucketModel",
		ui.Subtle("j/k, up/down: select")+" | "+ui.Subtle("enter: choose")+" | "+ui.Subtle("q, esc: quit"))
}
