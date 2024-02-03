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

var (
	currentNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle        = lipgloss.NewStyle().Margin(2, 1, 1)
	checkMark        = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
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
	// status is the status of the create bucket operation.
	status status
	// ctx is the context.
	ctx context.Context
	// err is the error that occurred during the operation.
	err error
	// width is the width of the terminal.
	window *ui.Window

	// TODO: refactor
	index    int
	sum      int
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
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		toggles:    ui.NewToggleSets(0),
		app:        app,
		ctx:        ctx,
		status:     statusNone,
		spinner:    s,
		progress:   p,
		index:      1,
		window:     ui.NewWindow(0, 0),
	}, nil
}

// Init initializes the model.
func (m *s3hubDeleteBucketModel) Init() tea.Cmd {
	return nil // Not called this method
}

// Update updates the model based on messages.
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
			m.status = statusQuit
			return m, tea.Quit
		case "q", "esc":
			m.status = statusReturnToTop
			return newRootModel(), nil
		case "enter":
			if m.status == statusReturnToTop || m.status == statusBucketDeleted {
				return newRootModel(), nil
			}

			if m.status == statusBucketListed {
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
				m.status = statusBucketDeleting
				m.index = 0 // Initialize index to 0 to accurately represent the starting state of progress.
				progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))

				return m, tea.Batch(
					m.spinner.Tick,
					progressCmd,
					tea.Printf("%s %s", checkMark, m.targetBuckets[0]),
					deleteS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
			}
		case " ":
			if m.status == statusBucketListed {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case tea.WindowSizeMsg:
		m.window.Width, m.window.Height = msg.Width, msg.Height
	case fetchS3BucketMsg:
		m.status = statusBucketFetched
		m.bucketSets = msg.buckets
		m.choice = ui.NewChoice(0, m.bucketSets.Len()-1)
		m.toggles = ui.NewToggleSets(m.bucketSets.Len())
		return m, nil
	case deleteS3BucketMsg:
		m.targetBuckets = m.targetBuckets[1:]
		if len(m.targetBuckets) == 0 {
			m.status = statusBucketDeleted
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
		m.status = statusQuit
		return m, tea.Quit
	default:
		return m, nil
	}
	return m, nil
}

func (m *s3hubDeleteBucketModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	switch m.status {
	case statusQuit:
		return ui.GoodByeMessage()
	case statusBucketDeleted:
		return doneStyle.Render("All S3 buckets deleted. Press <enter> to return to the top.\n")
	case statusBucketDeleting:
		w := lipgloss.Width(fmt.Sprintf("%d", m.sum))
		bucketCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.sum-1)

		spin := m.spinner.View() + " "
		prog := m.progress.View()
		cellsAvail := max(0, m.window.Width-lipgloss.Width(spin+prog+bucketCount))

		bucketName := currentNameStyle.Render(m.targetBuckets[0].String())
		info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Deleting " + bucketName)
		cellsRemaining := max(0, m.window.Width-lipgloss.Width(spin+info+prog+bucketCount))
		gap := strings.Repeat(" ", cellsRemaining)
		return spin + info + gap + prog + bucketCount
	case statusBucketFetching, statusNone:
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
	case statusBucketFetched:
		return m.bucketListString()
	default:
		return m.bucketListString() // TODO: implement
	}
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

	m.status = statusBucketListed
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
	m.status = statusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>, <esc>, q: return to the top"))
}
