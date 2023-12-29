package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
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
	template += subtle("j/k, up/down: select | enter: choose | q, <esc>: quit")

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

const (
	// s3hubCreateBucketRegionChoice is the choice number for selecting the AWS region.
	s3hubCreateBucketRegionChoice = 0
	// s3hubCreateBucketBucketNameChoice is the choice number for inputting the S3 bucket name.
	s3hubCreateBucketBucketNameChoice = 1
)

type s3hubCreateBucketModel struct {
	// bucketNameInput is the text input widget.
	bucketNameInput textinput.Model
	// err is the error that occurred during the operation.
	err error
	// bucket is the name of the S3 bucket that the user wants to create.
	bucket model.Bucket
	// state is the state of the create bucket operation.
	state s3hubCreateBucketState
	// awsConfig is the AWS configuration.
	awsConfig *model.AWSConfig
	// awsProfile is the AWS profile.
	awsProfile model.AWSProfile
	// region is the AWS region that the user wants to create the S3 bucket.
	region model.Region
	// choice is the currently selected menu item.
	choice int
	// app is the S3 application service.
	app *di.S3App
	ctx context.Context
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

	ctx := context.Background()
	profile := model.NewAWSProfile("")
	cfg, err := model.NewAWSConfig(ctx, profile, "")
	if err != nil {
		return nil, err
	}

	return &s3hubCreateBucketModel{
		bucketNameInput: ti,
		choice:          s3hubCreateBucketBucketNameChoice,
		awsConfig:       cfg,
		awsProfile:      profile,
		region:          cfg.Region(),
		ctx:             ctx,
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
		switch msg.String() {
		case "down":
			m.choice++
			if m.choice > s3hubCreateBucketBucketNameChoice {
				m.choice = s3hubCreateBucketRegionChoice
			}
		case "up":
			m.choice--
			if m.choice < s3hubCreateBucketRegionChoice {
				m.choice = s3hubCreateBucketBucketNameChoice
			}
		case "h", "left":
			if m.choice == s3hubCreateBucketRegionChoice {
				m.region = m.region.Prev()
			}
		case "l", "right":
			if m.choice == s3hubCreateBucketRegionChoice {
				m.region = m.region.Next()
			}
		case "enter":
			if m.bucketNameInput.Value() == "" || len(m.bucketNameInput.Value()) < model.BucketMinLength {
				return m, nil
			}

			app, err := di.NewS3App(m.ctx, m.awsProfile, m.region)
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.app = app
			m.bucket = model.Bucket(m.bucketNameInput.Value())
			return m, m.createS3BucketCmd()
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	case createMsg:
		m.state = s3hubCreateBucketStateCreated
		return m, tea.Quit
	}

	if m.choice == s3hubCreateBucketBucketNameChoice {
		var cmd tea.Cmd
		m.bucketNameInput, cmd = m.bucketNameInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *s3hubCreateBucketModel) View() string {
	if m.err != nil {
		message := fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   S3 Name   ]%s\n\n%s\n\n%s\n%s\n\n",
			m.awsProfile.String(),
			m.region.String(),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
			subtle("<enter>: create bucket"))

		message += fmt.Sprintf("%s\n", red("[Error]"))
		for _, line := range split(m.err.Error()) {
			message += fmt.Sprintf("  %s\n", red(line))
		}
		return message
	}

	if m.state == s3hubCreateBucketStateCreated {
		return fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   S3 Name   ]%s\n\n%s\n\n%s\n%s\n\n%s%s\n",
			m.awsProfile.String(),
			m.region.String(),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
			subtle("<enter>: create bucket"),
			"Created S3 bucket: ",
			yellow(m.bucket.String()))
	}

	if m.state == s3hubCreateBucketStateCreating {
		return fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   %s   ]%s\n\n%s\n\n%s\n%s\n\n%s\n",
			m.awsProfile.String(),
			m.region.String(),
			yellow("S3 Name"),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
			subtle("<enter>: create bucket"),
			"Creating S3 bucket...",
		)
	}

	if m.choice == s3hubCreateBucketRegionChoice {
		return fmt.Sprintf(
			"[ AWS Profile ] %s\n[ ◀︎  %s ▶︎ ] %s\n[   S3 Name   ]%s\n\n%s\n\n%s\n%s\n",
			m.awsProfile.String(),
			yellow("Region"),
			green(m.region.String()),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
			subtle("<enter>: create bucket | h/l, left/right: select region"),
		)
	}

	return fmt.Sprintf(
		"[ AWS Profile ] %s\n[    Region   ] %s\n[   %s   ]%s\n\n%s\n\n%s\n%s\n",
		m.awsProfile.String(),
		m.region.String(),
		yellow("S3 Name"),
		m.bucketNameWithColor(),
		m.bucketNameLengthString(),
		subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
		subtle("<enter>: create bucket"),
	)
}

// bucketNameWithColor returns the bucket name with color.
func (m *s3hubCreateBucketModel) bucketNameWithColor() string {
	if m.state == s3hubCreateBucketStateCreating || m.state == s3hubCreateBucketStateCreated {
		return m.bucketNameInput.View()
	}

	if len(m.bucketNameInput.Value()) < model.BucketMinLength && m.choice == s3hubCreateBucketBucketNameChoice {
		return red(m.bucketNameInput.View())
	}
	if m.choice == s3hubCreateBucketRegionChoice {
		return m.bucketNameInput.View()
	}
	return green(m.bucketNameInput.View())
}

// bucketNameLengthString returns the bucket name length string.
func (m *s3hubCreateBucketModel) bucketNameLengthString() string {
	lengthStr := fmt.Sprintf("Length: %d", len(m.bucketNameInput.Value()))
	if len(m.bucketNameInput.Value()) == model.BucketMaxLength {
		lengthStr += " (max)"
	} else if len(m.bucketNameInput.Value()) < model.BucketMinLength {
		lengthStr += " (min: 3)"
	}
	return lengthStr
}

func (m *s3hubCreateBucketModel) createS3BucketCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if m.app == nil {
			return errMsg(fmt.Errorf("not initialized s3 application. please restart the application"))
		}
		input := &usecase.S3BucketCreatorInput{
			Bucket: m.bucket,
			Region: m.region,
		}
		m.state = s3hubCreateBucketStateCreating

		if _, err := m.app.S3BucketCreator.CreateS3Bucket(m.ctx, input); err != nil {
			return errMsg(err)
		}
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
		subtle("j/k, up/down: select")+" | "+subtle("enter: choose")+" | "+subtle("q, esc: quit"))
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
		subtle("j/k, up/down: select")+" | "+subtle("enter: choose")+" | "+subtle("q, esc: quit"))
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
		subtle("j/k, up/down: select")+" | "+subtle("enter: choose")+" | "+subtle("q, esc: quit"))
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
		subtle("j/k, up/down: select")+" | "+subtle("enter: choose")+" | "+subtle("q, esc: quit"))

}
