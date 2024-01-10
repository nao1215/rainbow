package s3hub

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/ui"
)

type s3hubListBucketModel struct {
	// err is the error that occurred during the operation.
	err error
	// awsConfig is the AWS configuration.
	awsConfig *model.AWSConfig
	// awsProfile is the AWS profile.
	awsProfile model.AWSProfile
	// region is the AWS region that the user wants to create the S3 bucket.
	region model.Region
	// choice is the currently selected menu item.
	choice *ui.Choice
	// app is the S3 application service.
	app *di.S3App
	// ctx is the context.
	ctx context.Context
	// bucketSets is the list of the S3 buckets.
	bucketSets model.BucketSets
	// status is the status of the list bucket operation.
	status s3hubListBucketStatus
}

// s3hubListBucketStatus is the status of the list bucket operation.
type s3hubListBucketStatus int

// fetchMsg is the message that is sent when the user wants to fetch the list of the S3 buckets.
type fetchMsg struct{}

const (
	// s3hubListBucketStatusNone is the status when the list bucket operation is not executed.
	s3hubListBucketStatusNone s3hubListBucketStatus = iota
	// s3hubListBucketStatusBucketCreating is the status when the list bucket operation is executed and the bucket is being created.
	s3hubListBucketStatusBucketCreating
	// s3hubListBucketStatusBucketCreated is the status when the list bucket operation is executed and the bucket is created.
	s3hubListBucketStatusBucketCreated
	// s3hubListBucketStatusBucketListed is the status when the list bucket operation is executed and the bucket list is displayed.
	s3hubListBucketStatusBucketListed
	// s3hubListBucketStatusObjectListed is the status when the list bucket operation is executed and the object list is displayed.
	s3hubListBucketStatusObjectListed
	// s3hubListBucketStatusReturnToTop is the status when the user returns to the top.
	s3hubListBucketStatusReturnToTop
	// s3hubListBucketStatusQuit is the status when the user quits the application.
	s3hubListBucketStatusQuit
)

const (
	windowHeight = 10
)

func newS3HubListBucketModel() (*s3hubListBucketModel, error) {
	ctx := context.Background()
	profile := model.NewAWSProfile("")
	cfg, err := model.NewAWSConfig(ctx, profile, "")
	if err != nil {
		return nil, err
	}
	region := cfg.Region()

	app, err := di.NewS3App(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	return &s3hubListBucketModel{
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		app:        app,
		choice:     ui.NewChoice(0, 0),
		status:     s3hubListBucketStatusNone,
		ctx:        ctx,
		bucketSets: model.BucketSets{},
	}, nil
}

func (m *s3hubListBucketModel) Init() tea.Cmd {
	return nil // Not called this method
}

func (m *s3hubListBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice.Increment()
		case "k", "up":
			m.choice.Decrement()
		case "ctrl+c":
			m.status = s3hubListBucketStatusQuit
			return m, tea.Quit
		case "q", "esc":
			m.status = s3hubListBucketStatusReturnToTop
			return newRootModel(), nil
		case "enter":
			if m.status == s3hubListBucketStatusReturnToTop {
				return newRootModel(), nil
			}
		}
	case fetchMsg:
		return m, nil
	case ui.ErrMsg:
		m.err = msg
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

func (m *s3hubListBucketModel) View() string {
	if m.err != nil {
		m.status = s3hubListBucketStatusQuit
		return fmt.Sprintf("%s", ui.ErrorMessage(m.err))
	}

	if m.status == s3hubListBucketStatusQuit {
		return ui.GoodByeMessage()
	}

	if m.status == s3hubListBucketStatusNone || m.status == s3hubListBucketStatusBucketCreating {
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
	}

	if m.status == s3hubListBucketStatusBucketCreated {
		return m.bucketListString()
	}
	return m.bucketListString() // TODO: implement
}

// bucketListString returns the string representation of the bucket list.
func (m *s3hubListBucketModel) bucketListString() string {
	switch len(m.bucketSets) {
	case 0:
		return m.emptyBucketListString()
	default:
		return m.bucketListStrWithCheckbox()
	}
}

// bucketListStrWithCheckbox generates the string representation of the bucket list.
func (m *s3hubListBucketModel) bucketListStrWithCheckbox() string {
	startIndex := 0
	endIndex := len(m.bucketSets)

	if m.choice.Choice >= windowHeight {
		startIndex = m.choice.Choice - windowHeight + 1
		endIndex = startIndex + windowHeight
		if endIndex > len(m.bucketSets) {
			startIndex = len(m.bucketSets) - windowHeight
			endIndex = len(m.bucketSets)
		}
	} else {
		if len(m.bucketSets) > windowHeight {
			endIndex = windowHeight
		}
	}

	m.status = s3hubListBucketStatusBucketListed
	s := fmt.Sprintf("S3 buckets %d/%d (profile=%s)\n", m.choice.Choice+1, m.bucketSets.Len(), m.awsProfile.String())
	for i := startIndex; i < endIndex; i++ {
		b := m.bucketSets[i]
		s += fmt.Sprintf("%s\n",
			ui.Checkbox(
				fmt.Sprintf(
					"%s (region=%s, updated_at=%s)",
					color.GreenString("%s", b.Bucket),
					color.YellowString("%s", b.Region),
					b.CreationDate.Format("2006-01-02 15:04:05 MST")),
				m.choice.Choice == i))
	}
	s += ui.Subtle("\n<esc>: return to the top | <Ctrl-C>: quit | up/down: select\n")
	s += ui.Subtle("<enter>: choose bucket\n\n")
	return s
}

// emptyBucketListString returns the string representation when there are no S3 buckets.
func (m *s3hubListBucketModel) emptyBucketListString() string {
	m.status = s3hubListBucketStatusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}

// fetchS3BucketListCmd fetches the list of the S3 buckets.
func (m *s3hubListBucketModel) fetchS3BucketListCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.status = s3hubListBucketStatusBucketCreating

		output, err := m.app.S3BucketLister.ListS3Buckets(m.ctx, &usecase.S3BucketListerInput{})
		if err != nil {
			m.status = s3hubListBucketStatusQuit
			return ui.ErrMsg(err)
		}
		m.bucketSets = output.Buckets
		m.status = s3hubListBucketStatusBucketCreated
		m.choice = ui.NewChoice(0, len(m.bucketSets)-1)

		return fetchMsg{}
	})
}
