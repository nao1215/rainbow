package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/nao1215/rainbow/app/domain/model"
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
	// choice is the currently selected menu item.
	choice int
	// chosen is true when the user has chosen a menu item.
	chosen bool
	// quitting is true when the user has quit the application.
	quitting bool
	// err is the error that occurred during the operation.
	err error
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
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m.updateChoices(msg)
}

// View renders the application's UI.
func (m *s3hubRootModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s", m.err.Error())
	}

	if m.quitting {
		return "\n  See you later! (TODO: output log)\n\n" // TODO: print log.
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
			m.choice++
			if m.choice > s3hubTopMaxChoice {
				m.choice = s3hubTopMinChoice
			}
		case "k", "up":
			m.choice--
			if m.choice < s3hubTopMinChoice {
				m.choice = s3hubTopMaxChoice
			}
		case "enter":
			m.chosen = true
			switch m.choice {
			case s3hubTopCreateChoice:
				model, err := newS3hubCreateBucketModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				return model, nil
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

// choicesView returns a string containing the choices menu.
func (m *s3hubRootModel) choicesView() string {
	c := m.choice
	template := "%s\n\n"
	template += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, <esc>: quit")

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
	// textInput is the text input widget.
	textInput textinput.Model
	// err is the error that occurred during the operation.
	err error
	// bucketName is the name of the S3 bucket that the user wants to create.
	bucketName string
	// state is the state of the create bucket operation.
	state s3hubCreateBucketState
	// awsConfig is the AWS configuration.
	awsConfig *model.AWSConfig
	// awsProfile is the AWS profile.
	awsProfile model.AWSProfile
}

// createMsg is the message that is sent when the user wants to create the S3 bucket.
type createMsg struct{}

type s3hubCreateBucketState int

const (
	s3hubCreateBucketStateNone     s3hubCreateBucketState = 0
	s3hubCreateBucketStateCreating s3hubCreateBucketState = 1
	s3hubCreateBucketStateCreated  s3hubCreateBucketState = 2
)

func newS3hubCreateBucketModel() (*s3hubCreateBucketModel, error) {
	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("Write the S3 bucket name here (min: %d, max: %d)", model.BucketMinLength, model.BucketMaxLength)
	ti.Focus()
	ti.CharLimit = model.BucketMaxLength
	ti.Width = model.BucketMaxLength

	profile := model.NewAWSProfile("")
	cfg, err := model.NewAWSConfig(context.Background(), profile, "")
	if err != nil {
		return nil, err
	}

	return &s3hubCreateBucketModel{
		textInput:  ti,
		awsConfig:  cfg,
		awsProfile: profile,
	}, nil
}

func (m *s3hubCreateBucketModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *s3hubCreateBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() == "" || len(m.textInput.Value()) < model.BucketMinLength {
				return m, nil
			}
			m.bucketName = m.textInput.Value()
			m.state = s3hubCreateBucketStateCreating
			return m, createS3BucketCmd()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	case createMsg:
		// TODO: Wait for the result of the create bucket operation.
		m.state = s3hubCreateBucketStateCreated
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *s3hubCreateBucketModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("%s", m.err.Error())
	}

	if m.state == s3hubCreateBucketStateCreated {
		return fmt.Sprintf("Created S3 bucket: %s\n", m.bucketName)
	}

	if m.bucketName != "" {
		return fmt.Sprintf("Creating S3 bucket: %s (TODO: not implemented)\n", m.bucketName)
	}

	lengthStr := fmt.Sprintf("Length: %d", len(m.textInput.Value()))
	if len(m.textInput.Value()) == model.BucketMaxLength {
		lengthStr += " (max)"
	} else if len(m.textInput.Value()) < model.BucketMinLength {
		lengthStr += " (min: 3)"
	}

	return fmt.Sprintf(
		"[ AWS Profile ] %s\n[    Region   ] %s\n[Input S3 name] %s\n\n%s\n\n%s",
		m.awsProfile.String(), m.awsConfig.Region().String(),
		m.textInput.View(), lengthStr, subtle("<esc>, <Ctrl-C>: quit"),
	)
}

func createS3BucketCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// TODO: implement create s3 bucket operation.
		return createMsg{}
	})
}

type s3hubListBucketModel struct {
	// quitting is true when the user has quit the application.
	quitting bool
}

func (m *s3hubListBucketModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubListBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
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
	// quitting is true when the user has quit the application.
	quitting bool
}

func (m *s3hubCopyModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubCopyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
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
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))
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
		subtle("j/k, up/down: select")+dot+subtle("enter: choose")+dot+subtle("q, esc: quit"))

}
