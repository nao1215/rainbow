// Package s3hub is the text-based user interface for s3hub command.
package s3hub

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/nao1215/rainbow/ui"
)

const (
	// s3hubTopMinChoice is the minimum choice number.
	s3hubTopMinChoice = 0
	// s3hubTopMaxChoice is the maximum choice number.
	s3hubTopMaxChoice = 3
	// s3hubTopCreateChoice is the choice number for creating the S3 bucket.
	s3hubTopCreateChoice = 0
	// s3hubTopListChoice is the choice number for listing S3 buckets.
	s3hubTopListChoice = 1
	// s3hubTopDeleteContentsChoice is the choice number for deleting contents from the S3 bucket.
	s3hubTopDeleteContentsChoice = 2
	// s3hubTopDeleteBucketChoice is the choice number for deleting the S3 bucket.
	s3hubTopDeleteBucketChoice = 3
)

// s3hubRootModel is the top-level model for the application.
type s3hubRootModel struct {
	// choice is the currently selected menu item.
	choice *ui.Choice
	// chosen is true when the user has chosen a menu item.
	chosen bool
	// quitting is true when the user has quit the application.
	quitting bool
	// err is the error that occurred during the operation.
	err error
}

// RunS3hubUI start s3hub command interactive UI.
func RunS3hubUI() error {
	_, err := tea.NewProgram(newRootModel()).Run()
	return err
}

func newRootModel() *s3hubRootModel {
	return &s3hubRootModel{
		choice: ui.NewChoice(s3hubTopMinChoice, s3hubTopMaxChoice),
	}
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
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m.updateChoices(msg)
}

// View renders the application's UI.
func (m *s3hubRootModel) View() string {
	if m.err != nil {
		return ui.ErrorMessage(m.err)
	}

	if m.quitting {
		return ui.GoodByeMessage()
	}

	var s string
	if !m.chosen {
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
			m.choice.Increment()
		case "k", "up":
			m.choice.Decrement()
		case "enter":
			m.chosen = true
			switch m.choice.Choice {
			case s3hubTopCreateChoice:
				model, err := newS3hubCreateBucketModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				return model, nil
			case s3hubTopListChoice:
				model, err := newS3HubListBucketModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				model.s3BucketListBucketStatus = s3hubListBucketStatusBucketFetching
				return model, fetchS3BucketListCmd(model.ctx, model.app)
			case s3hubTopDeleteContentsChoice:
				return &s3hubDeleteContentsModel{}, nil
			case s3hubTopDeleteBucketChoice:
				model, err := newS3hubDeleteBucketModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				model.s3bucketListStatus = s3hubListBucketStatusBucketFetching
				return model, fetchS3BucketListCmd(model.ctx, model.app)
			}
		}
	}
	return m, nil
}

// choicesView returns a string containing the choices menu.
func (m *s3hubRootModel) choicesView() string {
	c := m.choice.Choice
	template := "%s\n"
	template += ui.Subtle("j/k, up/down: select | enter: choose | q, <esc>: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n",
		ui.Checkbox("Create the S3 bucket", c == s3hubTopMinChoice),
		ui.Checkbox("List and download S3 objects", c == 1),
		ui.Checkbox("Delete contents from the S3 bucket", c == 2),
		ui.Checkbox("Delete the S3 bucket", c == s3hubTopMaxChoice),
	)
	return fmt.Sprintf(template, choices)
}
