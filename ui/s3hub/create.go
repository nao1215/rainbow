package s3hub

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/ui"
)

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
	// state is the status of the create bucket operation.
	status status
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
	// ctx is the context.
	ctx context.Context
}

// newS3hubCreateBucketModel creates a new s3hubCreateBucketModel.
func newS3hubCreateBucketModel() (*s3hubCreateBucketModel, error) {
	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("Write the S3 bucket name here (min: %d, max: %d)", model.MinBucketNameLength, model.MaxBucketNameLength)
	ti.Focus()
	ti.CharLimit = model.MaxBucketNameLength
	ti.Width = model.MaxBucketNameLength

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

// Init initializes the model.
func (m *s3hubCreateBucketModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the model.
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
			if m.status == statusBucketCreated {
				return newRootModel(), nil
			}
			if m.bucketNameInput.Value() == "" || len(m.bucketNameInput.Value()) < model.MinBucketNameLength {
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
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return newRootModel(), nil
		}
	case ui.ErrMsg:
		m.err = msg
		return m, nil
	case createMsg:
		m.status = statusBucketCreated
		return m, nil
	}

	if m.status != statusBucketCreated && m.choice == s3hubCreateBucketBucketNameChoice {
		var cmd tea.Cmd
		m.bucketNameInput, cmd = m.bucketNameInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View renders the application's UI.
func (m *s3hubCreateBucketModel) View() string {
	if m.err != nil {
		message := fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   S3 Name   ]%s\n\n%s\n\n%s\n%s\n\n",
			m.awsProfile.String(),
			m.region.String(),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			ui.Subtle("<esc>: return to the top | <Ctrl-C>: quit | up/down: select"),
			ui.Subtle("<enter>: create bucket"))

		message += ui.ErrorMessage(m.err)
		return message
	}
	switch m.status {
	case statusBucketCreating:
		return fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   %s   ]%s\n\n%s\n\n%s\n%s\n\n%s\n",
			m.awsProfile.String(),
			m.region.String(),
			ui.Yellow("S3 Name"),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			ui.Subtle("<esc>, <Ctrl-C>: quit  | up/down: select"),
			ui.Subtle("<enter>: create bucket"),
			"Creating S3 bucket...",
		)
	case statusBucketCreated:
		return fmt.Sprintf("[ AWS Profile ] %s\n[    Region   ] %s\n[   S3 Name   ]%s\n\n%s\n\nCreated S3 bucket: %s\n%s\n",
			m.awsProfile.String(),
			m.region.String(),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			ui.Yellow(m.bucket.String()),
			ui.Subtle("<enter>: return to the top"))
	default:
		if m.choice == s3hubCreateBucketRegionChoice {
			return fmt.Sprintf(
				"[ AWS Profile ] %s\n[ ◀︎  %s ▶︎ ] %s\n[   S3 Name   ]%s\n\n%s\n\n%s\n%s\n",
				m.awsProfile.String(),
				ui.Yellow("Region"),
				ui.Green(m.region.String()),
				m.bucketNameWithColor(),
				m.bucketNameLengthString(),
				ui.Subtle("<esc>: return to the top | <Ctrl-C>: quit | up/down: select"),
				ui.Subtle("<enter>: create bucket   | h/l, left/right: select region"),
			)
		}
		return fmt.Sprintf(
			"[ AWS Profile ] %s\n[    Region   ] %s\n[   %s   ]%s\n\n%s\n\n%s\n%s\n",
			m.awsProfile.String(),
			m.region.String(),
			ui.Yellow("S3 Name"),
			m.bucketNameWithColor(),
			m.bucketNameLengthString(),
			ui.Subtle("<esc>: return to the top | <Ctrl-C>: quit | up/down: select"),
			ui.Subtle("<enter>: create bucket"),
		)
	}
}

// bucketNameWithColor returns the bucket name with color.
func (m *s3hubCreateBucketModel) bucketNameWithColor() string {
	if m.status == statusBucketCreating || m.status == statusBucketCreated {
		return m.bucketNameInput.View()
	}

	if len(m.bucketNameInput.Value()) < model.MinBucketNameLength && m.choice == s3hubCreateBucketBucketNameChoice {
		return ui.Red(m.bucketNameInput.View())
	}
	if m.choice == s3hubCreateBucketRegionChoice {
		return m.bucketNameInput.View()
	}
	return ui.Green(m.bucketNameInput.View())
}

// bucketNameLengthString returns the bucket name length string.
func (m *s3hubCreateBucketModel) bucketNameLengthString() string {
	lengthStr := fmt.Sprintf("Length: %d", len(m.bucketNameInput.Value()))
	if len(m.bucketNameInput.Value()) == model.MaxBucketNameLength {
		lengthStr += " (max)"
	} else if len(m.bucketNameInput.Value()) < model.MinBucketNameLength {
		lengthStr += " (min: 3)"
	}
	return lengthStr
}
