package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

const (
	// s3hubTopMinChoice is the minimum choice number.
	s3hubTopMinChoice = 0
	// s3hubTopMaxChoice is the maximum choice number.
	s3hubTopMaxChoice = 4
	// s3hubTopCreateChoice is the choice number for creating the S3 bucket.
	s3hubTopCreateChoice = 0
	// s3hubTopListChoice is the choice number for listing S3 buckets.
	s3hubTopListChoice = 1
	// s3hubTopCopyChoice is the choice number for copying file to the S3 bucket.
	s3hubTopCopyChoice = 2
	// s3hubTopDeleteContentsChoice is the choice number for deleting contents from the S3 bucket.
	s3hubTopDeleteContentsChoice = 3
	// s3hubTopDeleteBucketChoice is the choice number for deleting the S3 bucket.
	s3hubTopDeleteBucketChoice = 4
)

// s3hubRootModel is the top-level model for the application.
type s3hubRootModel struct {
	// Choice is the currently selected menu item.
	Choice int
	// Chosen is true when the user has chosen a menu item.
	Chosen bool
	// Quitting is true when the user has quit the application.
	Quitting bool
}

// RunS3hubUI start s3hub command interactive UI.
func RunS3hubUI() error {
	_, err := tea.NewProgram(&s3hubRootModel{}).Run()
	return err
}

// Init initializes the model.
func (m *s3hubRootModel) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m *s3hubRootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	if !m.Chosen {
		return m.updateChoices(msg)
	}
	return m.updateChosen(msg)
}

// View renders the application's UI.
func (m *s3hubRootModel) View() string {
	if m.Quitting {
		return "\n  See you later!\n\n" // TODO: print log.
	}

	var s string
	if !m.Chosen {
		s = m.choicesView()
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// updateChoices updates the model based on keypresses.
func (m *s3hubRootModel) updateChoices(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > s3hubTopMaxChoice {
				m.Choice = s3hubTopMinChoice
			}
		case "k", "up":
			m.Choice--
			if m.Choice < s3hubTopMinChoice {
				m.Choice = s3hubTopMaxChoice
			}
		case "enter":
			m.Chosen = true
			switch m.Choice {
			case s3hubTopCreateChoice:
				return &s3hubCreateBucketModel{}, nil
			case s3hubTopListChoice:
				return &s3hubListBucketModel{}, nil
			case s3hubTopCopyChoice:
				return &s3hubCopyModel{}, nil
			case s3hubTopDeleteContentsChoice:
				return &s3hubDeleteContentsModel{}, nil
			case s3hubTopDeleteBucketChoice:
				return &s3hubDeleteBucketModel{}, nil
			}
		}
	}
	return m, nil
}

// updateChosen updates the model when the user has chosen a menu item.
func (m *s3hubRootModel) updateChosen(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	}
	return m, nil
}

// choicesView returns a string containing the choices menu.
func (m *s3hubRootModel) choicesView() string {
	c := m.Choice
	template := "%s\n\n"
	template += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		checkbox("Create the S3 bucket", c == s3hubTopMinChoice),
		checkbox("List S3 buckets", c == 1),
		checkbox("Copy file to the S3 bucket", c == 2),
		checkbox("Delete contents from the S3 bucket", c == 3),
		checkbox("Delete the S3 bucket", c == s3hubTopMaxChoice),
	)
	return fmt.Sprintf(template, choices)
}

type s3hubCreateBucketModel struct {
	// Quitting is true when the user has quit the application.
	Quitting bool
}

func (m *s3hubCreateBucketModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubCreateBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		m.Quitting = true
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubCreateBucketModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubCreateBucketModel",
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))
}

type s3hubListBucketModel struct {
	// Quitting is true when the user has quit the application.
	Quitting bool
}

func (m *s3hubListBucketModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubListBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubListBucketModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubListBucketModel",
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))
}

type s3hubCopyModel struct {
	// Quitting is true when the user has quit the application.
	Quitting bool
}

func (m *s3hubCopyModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubCopyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubCopyModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubCopyModel",
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))
}

type s3hubDeleteContentsModel struct {
	// Quitting is true when the user has quit the application.
	Quitting bool
}

func (m *s3hubDeleteContentsModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubDeleteContentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubDeleteContentsModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubDeleteContentsModel",
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))
}

type s3hubDeleteBucketModel struct {
	// Quitting is true when the user has quit the application.
	Quitting bool
}

func (m *s3hubDeleteBucketModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubDeleteBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *s3hubDeleteBucketModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubDeleteBucketModel",
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))

}
