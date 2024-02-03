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

// s3hubListBucketStatus is the status of the list bucket operation.
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
	status status
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
}

const (
	windowHeight = 10
)

// newS3HubListBucketModel returns a new s3hubListBucketModel.
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
		status:     statusNone,
		ctx:        ctx,
		bucketSets: model.BucketSets{},
		toggles:    ui.NewToggleSets(0),
	}, nil
}

// Init initializes the model.
func (m *s3hubListBucketModel) Init() tea.Cmd {
	return nil // Not called this method
}

// Update updates the model.
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
			m.status = statusQuit
			return m, tea.Quit
		case "q", "esc":
			m.status = statusReturnToTop
			return newRootModel(), nil
		case "d":
			if m.status == statusBucketListed {
				m.status = statusDownloading

				buckets := make([]model.Bucket, 0, len(m.bucketSets))
				for i, b := range m.bucketSets {
					if m.toggles[i].Enabled {
						buckets = append(buckets, b.Bucket)
					}
				}
				return m, downloadS3BucketCmd(m.ctx, m.app, buckets)
			}
		case "enter":
			if m.status == statusReturnToTop || m.status == statusDownloaded {
				return newRootModel(), nil
			}
			if m.status == statusBucketListed {
				model, err := newS3HubListS3ObjectModel()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				model.status = statusS3ObjectFetching
				model.bucket = m.bucketSets[m.choice.Choice].Bucket
				return model, fetchS3KeysCmd(m.ctx, m.app, model.bucket)
			}
		case " ":
			if m.status == statusBucketListed {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case fetchS3BucketMsg:
		m.status = statusBucketFetched
		m.bucketSets = msg.buckets
		m.choice = ui.NewChoice(0, m.bucketSets.Len()-1)
		m.toggles = ui.NewToggleSets(m.bucketSets.Len())
		return m, nil
	case downloadS3BucketMsg:
		m.status = statusDownloaded
		return m, nil
	case ui.ErrMsg:
		m.err = msg
		m.status = statusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

// View renders the application's UI.
func (m *s3hubListBucketModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	if m.status == statusQuit {
		return ui.GoodByeMessage()
	}

	if m.status == statusDownloaded {
		return doneStyle.Render("All S3 buckets downloaded. Press <enter> to return to the top.")
	}

	if m.status == statusNone || m.status == statusBucketFetching {
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
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

	m.status = statusBucketListed
	s := fmt.Sprintf("S3 buckets %d/%d (profile=%s)\n\n", m.choice.Choice+1, m.bucketSets.Len(), m.awsProfile.String())
	for i := startIndex; i < endIndex; i++ {
		b := m.bucketSets[i]
		s += fmt.Sprintf("%s\n",
			ui.ToggleWidget(
				fmt.Sprintf(
					"%s (region=%s, updated_at=%s)",
					color.GreenString("%s", b.Bucket),
					color.YellowString("%s", b.Region),
					b.CreationDate.Format("2006-01-02 15:04:05 MST")),
				m.choice.Choice == i, m.toggles[i].Enabled))
	}
	s += ui.Subtle("\n<esc>: return to the top | <Ctrl-C>: quit | up/down: select\n")
	s += ui.Subtle("<space>: choose bucket to download | d: download buckets\n")
	s += ui.Subtle("<enter>: list up s3 objects in bucket\n\n")
	return s
}

// emptyBucketListString returns the string representation when there are no S3 buckets.
func (m *s3hubListBucketModel) emptyBucketListString() string {
	m.status = statusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}

type s3hubListS3ObjectModel struct {
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
	// bucket is the S3 bucket that the user wants to list the objects.
	bucket model.Bucket
	// s3Keys is the list of the S3 bucket objects.
	s3Keys []model.S3Key
	// status is the status of the list S3 object operation.
	status status
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
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
		toggles:    ui.NewToggleSets(0),
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
			model.status = statusBucketFetching
			return model, fetchS3BucketListCmd(model.ctx, model.app)
		case "d":
			if m.status == statusS3ObjectListed {
				m.status = statusDownloading
				keys := make([]model.S3Key, 0, len(m.s3Keys))
				for i, k := range m.s3Keys {
					if m.toggles[i].Enabled {
						keys = append(keys, k)
					}
				}
				return m, downloadS3ObjectsCmd(m.ctx, m.app, m.bucket, keys)
			}
		case "enter":
			if m.status == statusReturnToTop || m.status == statusDownloaded {
				return newRootModel(), nil
			}
		case " ":
			if m.status == statusS3ObjectListed {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case fetchS3Keys:
		m.status = statusS3ObjectFetched
		m.s3Keys = msg.keys
		m.choice = ui.NewChoice(0, len(m.s3Keys)-1)
		m.toggles = ui.NewToggleSets(len(m.s3Keys))
		return m, nil
	case downloadS3BucketMsg:
		m.status = statusDownloaded
		return m, nil
	case ui.ErrMsg:
		m.err = msg
		m.status = statusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

// View renders the application's UI.
func (m *s3hubListS3ObjectModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	if m.status == statusQuit {
		return ui.GoodByeMessage()
	}

	if m.status == statusDownloaded {
		return doneStyle.Render("All S3 objects downloaded. Press <enter> to return to the top.")
	}

	if m.status == statusNone || m.status == statusS3ObjectFetching {
		return fmt.Sprintf(
			"fetching the list of the S3 objects (profile=%s, bucket=%s)\n",
			m.awsProfile.String(),
			m.bucket.String())
	}
	if m.status == statusS3ObjectFetched {
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
	} else if len(m.s3Keys) > windowHeight {
		endIndex = windowHeight
	}

	m.status = statusS3ObjectListed
	s := fmt.Sprintf("S3 objects %d/%d (profile=%s)\n\n", m.choice.Choice+1, len(m.s3Keys), m.awsProfile.String())
	for i := startIndex; i < endIndex; i++ {
		s += fmt.Sprintf("%s\n",
			ui.ToggleWidget(color.GreenString("%s", m.bucket.Join(m.s3Keys[i])), m.choice.Choice == i, m.toggles[i].Enabled))
	}
	s += ui.Subtle("\n<esc>: return | <Ctrl-C>: quit | up/down: select\n")
	s += ui.Subtle("<space>: choose s3 object to download | d: download s3 object\n\n")
	return s
}

// emptyS3ObjectListString returns the string representation when there are no S3 objects.
func (m *s3hubListS3ObjectModel) emptyS3ObjectListString() string {
	m.status = statusReturnToTop
	return fmt.Sprintf("No S3 objects (profile=%s, bucket=%s)\n\n%s\n",
		m.awsProfile.String(),
		m.bucket.String(),
		ui.Subtle("<enter>, <esc>, q: return to the top"))
}
