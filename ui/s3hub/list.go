package s3hub

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
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

const (
	// s3hubListBucketStatusNone is the status when the list bucket operation is not executed.
	s3hubListBucketStatusNone s3hubListBucketStatus = iota
	// s3hubListBucketStatusBucketFetching is the status when the list bucket operation is executed.
	s3hubListBucketStatusBucketFetching
	// s3hubListBucketStatusBucketFetched is the status when the list bucket operation is executed and the bucket list is fetched.
	s3hubListBucketStatusBucketFetched
	// s3hubListBucketStatusBucketListed is the status when the list bucket operation is executed and the bucket list is displayed.
	s3hubListBucketStatusBucketListed
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
			if m.status == s3hubListBucketStatusBucketListed {
				model, err := newS3HubListS3ObjectModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				model.status = s3hubListS3ObjectStatusFetching
				model.bucket = m.bucketSets[m.choice.Choice].Bucket
				return model, fetchS3KeysCmd(m.ctx, m.app, model.bucket)
			}
		case "space":
			// TODO: implement
		}
	case fetchS3BucketMsg:
		m.status = s3hubListBucketStatusBucketFetched
		m.bucketSets = msg.buckets
		m.choice = ui.NewChoice(0, m.bucketSets.Len()-1)
		return m, nil
	case ui.ErrMsg:
		m.err = msg
		m.status = s3hubListBucketStatusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

func (m *s3hubListBucketModel) View() string {
	if m.err != nil {
		m.status = s3hubListBucketStatusQuit
		return ui.ErrorMessage(m.err)
	}

	if m.status == s3hubListBucketStatusQuit {
		return ui.GoodByeMessage()
	}

	if m.status == s3hubListBucketStatusNone || m.status == s3hubListBucketStatusBucketFetching {
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
	}

	if m.status == s3hubListBucketStatusBucketFetched {
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
	s := fmt.Sprintf("S3 buckets %d/%d (profile=%s)\n\n", m.choice.Choice+1, m.bucketSets.Len(), m.awsProfile.String())
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
	s += ui.Subtle("<enter>, <space>: choose bucket\n\n")
	return s
}

// emptyBucketListString returns the string representation when there are no S3 buckets.
func (m *s3hubListBucketModel) emptyBucketListString() string {
	m.status = s3hubListBucketStatusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}

// s3hubListS3ObjectStatus is the status of the list s3 objects operation.
type s3hubListS3ObjectStatus int

const (
	// s3hubListBucketStatusNone is the status when the list bucket operation is not executed.
	s3hubListS3ObjectStatusNone s3hubListS3ObjectStatus = iota
	// s3hubListS3ObjectStatusFetching is the status when the list bucket operation is executed.
	s3hubListS3ObjectStatusFetching
	// s3hubListS3ObjectStatusFetched is the status when the list bucket operation is executed and the bucket list is fetched.
	s3hubListS3ObjectStatusFetched
	// s3hubListBucketStatusBucketListed is the status when the list bucket operation is executed and the bucket list is displayed.
	s3hubListS3ObjectStatusListed
	// s3hubListBucketStatusReturnToTop is the status when the user returns to the top.
	s3hubListS3ObjectStatusReturnToTop
	// s3hubListBucketStatusQuit is the status when the user quits the application.
	s3hubListS3ObjectStatusQuit
)

type s3hubListS3ObjectModel struct {
	// err is the error that occurred during the operation.
	err error
	// status is the status of the list s3 objects operation.
	status s3hubListS3ObjectStatus
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
	// bucket is the S3 bucket that the user wants to list the objects.
	bucket model.Bucket
	// s3Keys is the list of the S3 bucket objects.
	s3Keys []model.S3Key
}

// newS3HubListS3ObjectModel returns a new s3hubListS3ObjectModel.
func newS3HubListS3ObjectModel() (*s3hubListS3ObjectModel, error) {
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

	return &s3hubListS3ObjectModel{
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		app:        app,
		choice:     ui.NewChoice(0, 0),
		ctx:        ctx,
	}, nil
}

// Init initializes the model.
func (m *s3hubListS3ObjectModel) Init() tea.Cmd {
	return nil // Not called this method
}

// Update updates the model.
func (m *s3hubListS3ObjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Quit
		case "q", "esc":
			model, err := newS3HubListBucketModel()
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			model.status = s3hubListBucketStatusBucketFetching
			return model, fetchS3BucketListCmd(model.ctx, model.app)
		}

	case fetchS3Keys:
		m.status = s3hubListS3ObjectStatusFetched
		m.s3Keys = msg.keys
		m.choice = ui.NewChoice(0, len(m.s3Keys)-1)
		return m, nil
	case ui.ErrMsg:
		m.err = msg
		m.status = s3hubListS3ObjectStatusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

// View renders the application's UI.
func (m *s3hubListS3ObjectModel) View() string {
	if m.err != nil {
		m.status = s3hubListS3ObjectStatusQuit
		return ui.ErrorMessage(m.err)
	}

	if m.status == s3hubListS3ObjectStatusQuit {
		return ui.GoodByeMessage()
	}

	if m.status == s3hubListS3ObjectStatusNone || m.status == s3hubListS3ObjectStatusFetching {
		return fmt.Sprintf(
			"fetching the list of the S3 objects (profile=%s, bucket=%s)\n",
			m.awsProfile.String(),
			m.bucket.String())
	}

	if m.status == s3hubListS3ObjectStatusFetched {
		return m.s3ObjectListString()
	}
	return m.s3ObjectListString()
}

// s3ObjectListString returns the string representation of the S3 object list.
func (m *s3hubListS3ObjectModel) s3ObjectListString() string {
	switch len(m.s3Keys) {
	case 0:
		return m.emptyS3ObjectListString()
	default:
		return m.s3ObjectListStrWithCheckbox()
	}
}

// s3ObjectListStrWithCheckbox generates the string representation of the S3 object list.
func (m *s3hubListS3ObjectModel) s3ObjectListStrWithCheckbox() string {
	startIndex := 0
	endIndex := len(m.s3Keys)

	if m.choice.Choice >= windowHeight {
		startIndex = m.choice.Choice - windowHeight + 1
		endIndex = startIndex + windowHeight
		if endIndex > len(m.s3Keys) {
			startIndex = len(m.s3Keys) - windowHeight
			endIndex = len(m.s3Keys)
		}
	} else {
		if len(m.s3Keys) > windowHeight {
			endIndex = windowHeight
		}
	}

	m.status = s3hubListS3ObjectStatusListed
	s := fmt.Sprintf("S3 objects %d/%d (profile=%s)\n\n", m.choice.Choice+1, len(m.s3Keys), m.awsProfile.String())
	for i := startIndex; i < endIndex; i++ {
		s += fmt.Sprintf("%s\n",
			ui.Checkbox(
				fmt.Sprintf(
					"%s",
					color.GreenString("%s", m.bucket.Join(m.s3Keys[i]))),
				m.choice.Choice == i))
	}
	s += ui.Subtle("\n<esc>: return | <Ctrl-C>: quit | up/down: select\n")
	s += ui.Subtle("<enter>, <space>: choose bucket\n\n")
	return s
}

// emptyS3ObjectListString returns the string representation when there are no S3 objects.
func (m *s3hubListS3ObjectModel) emptyS3ObjectListString() string {
	m.status = s3hubListS3ObjectStatusReturnToTop
	return fmt.Sprintf("No S3 objects (profile=%s, bucket=%s)\n\n%s\n",
		m.awsProfile.String(),
		m.bucket.String(),
		ui.Subtle("<esc>, q: return"))
}
