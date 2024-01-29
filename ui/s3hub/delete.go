package s3hub

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/ui"
)

type s3hubDeleteContentsModel struct {
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
}

// s3hubDeleteContentsStatus is the status of the delete contents operation.
func newS3hubDeleteContentsModel() (*s3hubDeleteContentsModel, error) {
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

	return &s3hubDeleteContentsModel{
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		app:        app,
		ctx:        ctx,
	}, nil
}

func (m *s3hubDeleteContentsModel) Init() tea.Cmd {
	return nil
}

func (m *s3hubDeleteContentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *s3hubDeleteContentsModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		"s3hubDeleteContentsModel",
		ui.Subtle("j/k, up/down: select")+" | "+ui.Subtle("enter: choose")+" | "+ui.Subtle("q, esc: quit"))
}

// s3hubDeleteBucketStatus is the status of the delete bucket operation.
type s3hubDeleteBucketStatus int

const (
	// s3hubDeleteBucketStatusNone is the status when the delete bucket operation is not executed.
	s3hubDeleteBucketStatusNone s3hubDeleteBucketStatus = iota
	// s3hubDeleteBucketStatusBucketDeleting is the status when the delete bucket operation is executed and the bucket is being deleted.
	s3hubDeleteBucketStatusBucketDeleting
	// s3hubDeleteBucketStatusBucketDeleted is the status when the delete bucket operation is executed and the bucket is deleted.
	s3hubDeleteBucketStatusBucketDeleted
)

var (
	currentBucketNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle              = lipgloss.NewStyle().Margin(2, 1, 1)
	checkMark              = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

type s3hubDeleteBucketModel struct {
	// awsConfig is the AWS configuration.
	awsConfig *model.AWSConfig
	// awsProfile is the AWS profile.
	awsProfile model.AWSProfile
	// region is the AWS region that the user wants to create the S3 bucket.
	region model.Region
	// choice is the currently selected menu item.
	choice *ui.Choice
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
	// app is the S3 application service.
	app *di.S3App
	// bucketSets is the list of the S3 buckets.
	bucketSets model.BucketSets
	// targetBuckets is the list of the S3 buckets that the user wants to delete.
	targetBuckets []model.Bucket
	// s3bucketListStatus is the status of the list bucket operation.
	s3bucketListStatus s3hubListBucketStatus
	// s3bucketDeleteStatus is the status of the delete bucket operation.
	s3bucketDeleteStatus s3hubDeleteBucketStatus
	// ctx is the context.
	ctx context.Context
	// err is the error that occurred during the operation.
	err error

	// TODO: refactor
	index    int
	sum      int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
}

// s3hubDeleteBucketStatus is the status of the delete bucket operation.
func newS3hubDeleteBucketModel() (*s3hubDeleteBucketModel, error) {
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

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return &s3hubDeleteBucketModel{
		awsConfig:          cfg,
		awsProfile:         profile,
		region:             region,
		toggles:            ui.NewToggleSets(0),
		app:                app,
		ctx:                ctx,
		s3bucketListStatus: s3hubListBucketStatusNone,
		spinner:            s,
		progress:           p,
		index:              1,
	}, nil
}

func (m *s3hubDeleteBucketModel) Init() tea.Cmd {
	return nil // Not called this method
}

func (m *s3hubDeleteBucketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.s3bucketListStatus = s3hubListBucketStatusQuit
			return m, tea.Quit
		case "q", "esc":
			m.s3bucketListStatus = s3hubListBucketStatusReturnToTop
			return newRootModel(), nil
		case "enter":
			if m.s3bucketListStatus == s3hubListBucketStatusReturnToTop || m.s3bucketDeleteStatus == s3hubDeleteBucketStatusBucketDeleted {
				return newRootModel(), nil
			}

			if m.s3bucketListStatus == s3hubListBucketStatusBucketListed && m.s3bucketDeleteStatus == s3hubDeleteBucketStatusNone {
				m.targetBuckets = make([]model.Bucket, 0, len(m.toggles))
				for i, t := range m.toggles {
					if t.Enabled {
						m.targetBuckets = append(m.targetBuckets, m.bucketSets[i].Bucket)
					}
				}
				if len(m.targetBuckets) == 0 {
					return m, nil
				}
				m.sum = len(m.targetBuckets) + 1
				m.s3bucketDeleteStatus = s3hubDeleteBucketStatusBucketDeleting
				return m, tea.Batch(m.spinner.Tick, deleteS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
			}
		case " ":
			if m.s3bucketListStatus == s3hubListBucketStatusBucketListed && m.s3bucketDeleteStatus == s3hubDeleteBucketStatusNone {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case fetchS3BucketMsg:
		m.s3bucketListStatus = s3hubListBucketStatusBucketFetched
		m.bucketSets = msg.buckets
		m.choice = ui.NewChoice(0, m.bucketSets.Len()-1)
		m.toggles = ui.NewToggleSets(m.bucketSets.Len())
		return m, nil
	case deleteS3BucketMsg:
		m.targetBuckets = m.targetBuckets[1:]
		if len(m.targetBuckets) == 0 {
			m.s3bucketDeleteStatus = s3hubDeleteBucketStatusBucketDeleted
			return m, nil
		}
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
		m.index++
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.targetBuckets[0]),
			deleteS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	case ui.ErrMsg:
		m.err = msg
		m.s3bucketListStatus = s3hubListBucketStatusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

func (m *s3hubDeleteBucketModel) View() string {
	if m.err != nil {
		m.s3bucketListStatus = s3hubListBucketStatusQuit
		return ui.ErrorMessage(m.err)
	}

	if m.s3bucketListStatus == s3hubListBucketStatusQuit {
		return ui.GoodByeMessage()
	}

	if m.s3bucketDeleteStatus == s3hubDeleteBucketStatusBucketDeleted {
		return doneStyle.Render("All S3 buckets deleted. Press <enter> to return to the top.\n")
	}

	if m.s3bucketDeleteStatus == s3hubDeleteBucketStatusBucketDeleting {
		w := lipgloss.Width(fmt.Sprintf("%d", m.sum))
		bucketCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.sum-1)

		spin := m.spinner.View() + " "
		prog := m.progress.View()
		cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+bucketCount))

		bucketName := currentBucketNameStyle.Render(m.targetBuckets[0].String())
		info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Deleting " + bucketName)
		cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+bucketCount))
		gap := strings.Repeat(" ", cellsRemaining)
		return spin + info + gap + prog + bucketCount
	}

	if m.s3bucketListStatus == s3hubListBucketStatusNone || m.s3bucketListStatus == s3hubListBucketStatusBucketFetching {
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
	}

	if m.s3bucketListStatus == s3hubListBucketStatusBucketFetched {
		return m.bucketListString()
	}
	return m.bucketListString() // TODO: implement
}

// bucketListString returns the string representation of the bucket list.
func (m *s3hubDeleteBucketModel) bucketListString() string {
	switch len(m.bucketSets) {
	case 0:
		return m.emptyBucketListString()
	default:
		return m.bucketListStrWithCheckbox()
	}
}

// bucketListStrWithCheckbox generates the string representation of the bucket list.
func (m *s3hubDeleteBucketModel) bucketListStrWithCheckbox() string {
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

	m.s3bucketListStatus = s3hubListBucketStatusBucketListed
	s := fmt.Sprintf("Select the S3 bucket(s) you want to delete %d/%d (profile=%s)\n\n",
		m.choice.Choice+1, m.bucketSets.Len(), m.awsProfile.String())
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
	s += ui.Subtle("<space>: choose bucket   | <enter> delete choosed bucket\n\n")
	return s
}

// emptyBucketListString returns the string representation when there are no S3 buckets.
func (m *s3hubDeleteBucketModel) emptyBucketListString() string {
	m.s3bucketListStatus = s3hubListBucketStatusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}
