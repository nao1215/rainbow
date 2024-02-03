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
	// targetBuckets is the list of the S3 buckets that the user wants to delete or download.
	targetBuckets []model.Bucket
	// status is the status of the list bucket operation.
	status status
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
	// width is the width of the terminal.
	window *ui.Window

	// TODO: refactor
	index    int
	sum      int
	spinner  spinner.Model
	progress progress.Model
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

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

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
		spinner:    s,
		progress:   p,
		index:      1,
		window:     ui.NewWindow(0, 0),
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
				m.targetBuckets = m.getTargetBuckets()
				if len(m.targetBuckets) == 0 {
					return m, nil
				}
				m.sum = len(m.targetBuckets) + 1
				m.status = statusDownloading
				progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))

				return m, tea.Batch(
					m.spinner.Tick,
					progressCmd,
					tea.Printf("%s %s", checkMark, m.targetBuckets[0]),
					downloadS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
			}
		case "D":
			if m.status == statusBucketListed {
				m.targetBuckets = m.getTargetBuckets()
				if len(m.targetBuckets) == 0 {
					return m, nil
				}
				m.sum = len(m.targetBuckets) + 1
				m.status = statusBucketDeleting
				progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))

				return m, tea.Batch(m.spinner.Tick,
					progressCmd,
					tea.Printf("%s %s", checkMark, m.targetBuckets[0]),
					deleteS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
			}
		case "enter":
			if m.status == statusReturnToTop || m.status == statusDownloaded || m.status == statusBucketDeleted {
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
	case tea.WindowSizeMsg:
		m.window.Width, m.window.Height = msg.Width, msg.Height
	case fetchS3BucketMsg:
		m.status = statusBucketFetched
		m.bucketSets = msg.buckets
		m.choice = ui.NewChoice(0, m.bucketSets.Len()-1)
		m.toggles = ui.NewToggleSets(m.bucketSets.Len())
		return m, nil
	case downloadS3BucketMsg:
		m.targetBuckets = m.targetBuckets[1:]
		if len(m.targetBuckets) == 0 {
			m.status = statusDownloaded
			return m, nil
		}
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
		m.index++
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.targetBuckets[0]),
			downloadS3BucketCmd(m.ctx, m.app, m.targetBuckets[0]))
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

// View renders the application's UI.
func (m *s3hubListBucketModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	switch m.status {
	case statusQuit:
		return ui.GoodByeMessage()
	case statusDownloaded:
		return doneStyle.Render("All S3 buckets downloaded. Press <enter> to return to the top.")
	case statusDownloading:
		w := lipgloss.Width(fmt.Sprintf("%d", m.sum))
		bucketCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.sum-1)
		spin := m.spinner.View() + " "
		prog := m.progress.View()
		cellsAvail := max(0, m.window.Width-lipgloss.Width(spin+prog+bucketCount))

		bucketName := currentNameStyle.Render(m.targetBuckets[0].String())
		info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Downloading " + bucketName)
		cellsRemaining := max(0, m.window.Width-lipgloss.Width(spin+info+prog+bucketCount))
		gap := strings.Repeat(" ", cellsRemaining)
		return spin + info + gap + prog + bucketCount
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
	case statusNone, statusBucketFetching:
		return fmt.Sprintf(
			"fetching the list of the S3 buckets (profile=%s)\n",
			m.awsProfile.String())
	case statusBucketFetched:
		return m.bucketListString()
	default:
		return m.bucketListString()
	}
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
	s += ui.Subtle("<space>: choose bucket to download or delete | <enter>: list up s3 objects in bucket\n")
	s += ui.Subtle("d: download buckets      | D: delete buckets\n\n")
	return s
}

// emptyBucketListString returns the string representation when there are no S3 buckets.
func (m *s3hubListBucketModel) emptyBucketListString() string {
	m.status = statusReturnToTop
	return fmt.Sprintf("No S3 buckets (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}

// getTargetBuckets returns the list of the S3 buckets that the user wants to delete or download
func (m *s3hubListBucketModel) getTargetBuckets() []model.Bucket {
	targetBuckets := make([]model.Bucket, 0, len(m.toggles))
	for i, t := range m.toggles {
		if t.Enabled {
			targetBuckets = append(targetBuckets, m.bucketSets[i].Bucket)
		}
	}
	return targetBuckets
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
	// targetS3Keys is the list of the S3 bucket objects that the user wants to download.
	targetS3Keys []model.S3Key
	// status is the status of the list S3 object operation.
	status status
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
	// width is the width of the terminal.
	window *ui.Window

	// TODO: refactor
	index    int
	sum      int
	spinner  spinner.Model
	progress progress.Model
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

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return &s3hubListS3ObjectModel{
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		app:        app,
		choice:     ui.NewChoice(0, 0),
		ctx:        ctx,
		toggles:    ui.NewToggleSets(0),
		spinner:    s,
		progress:   p,
		index:      1,
		window:     ui.NewWindow(0, 0),
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
				m.targetS3Keys = m.getTargetS3Keys()
				if len(m.targetS3Keys) == 0 {
					return m, nil
				}
				m.sum = len(m.targetS3Keys) + 1
				m.status = statusDownloading

				progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
				return m, tea.Batch(m.spinner.Tick,
					progressCmd,
					tea.Printf("%s %s", checkMark, m.targetS3Keys[0]),
					downloadS3ObjectsCmd(m.ctx, m.app, m.bucket, m.targetS3Keys[0]))
			}
		case "D":
			if m.status == statusS3ObjectListed {
				m.targetS3Keys = m.getTargetS3Keys()
				if len(m.targetS3Keys) == 0 {
					return m, nil
				}
				m.sum = len(m.targetS3Keys) + 1
				m.status = statusS3ObjectDeleting

				progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
				return m, tea.Batch(m.spinner.Tick,
					progressCmd,
					tea.Printf("%s %s", checkMark, m.targetS3Keys[0]),
					deleteS3ObjectCmd(m.ctx, m.app, m.bucket, m.targetS3Keys[0]))
			}
		case "enter":
			if m.status == statusReturnToTop || m.status == statusDownloaded || m.status == statusS3ObjectDeleted {
				return newRootModel(), nil
			}
		case " ":
			if m.status == statusS3ObjectListed {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case tea.WindowSizeMsg:
		m.window.Width, m.window.Height = msg.Width, msg.Height
	case fetchS3Keys:
		m.status = statusS3ObjectFetched
		m.s3Keys = msg.keys
		m.choice = ui.NewChoice(0, len(m.s3Keys)-1)
		m.toggles = ui.NewToggleSets(len(m.s3Keys))
		return m, nil
	case downloadS3ObjectsMsg:
		m.targetS3Keys = m.targetS3Keys[1:]
		if len(m.targetS3Keys) == 0 {
			m.status = statusDownloaded
			return m, nil
		}
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
		m.index++
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.targetS3Keys[0]),
			downloadS3ObjectsCmd(m.ctx, m.app, m.bucket, m.targetS3Keys[0]))
	case deleteS3ObjectMsg:
		m.targetS3Keys = m.targetS3Keys[1:]
		if len(m.targetS3Keys) == 0 {
			m.status = statusS3ObjectDeleted
			return m, nil
		}
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(m.sum-1))
		m.index++
		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.targetS3Keys[0]),
			deleteS3ObjectCmd(m.ctx, m.app, m.bucket, m.targetS3Keys[0]))
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

// View renders the application's UI.
func (m *s3hubListS3ObjectModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	switch m.status {
	case statusQuit:
		return ui.GoodByeMessage()
	case statusDownloaded:
		return doneStyle.Render("All S3 objects downloaded. Press <enter> to return to the top.")
	case statusDownloading:
		w := lipgloss.Width(fmt.Sprintf("%d", m.sum))
		s3keyCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.sum-1)
		spin := m.spinner.View() + " "
		prog := m.progress.View()
		cellsAvail := max(0, m.window.Width-lipgloss.Width(spin+prog+s3keyCount))

		s3keyName := currentNameStyle.Render(m.targetS3Keys[0].String())
		info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Downloading " + s3keyName)
		cellsRemaining := max(0, m.window.Width-lipgloss.Width(spin+info+prog+s3keyCount))
		gap := strings.Repeat(" ", cellsRemaining)
		return spin + info + gap + prog + s3keyCount
	case statusS3ObjectDeleted:
		return doneStyle.Render("All S3 objects deleted. Press <enter> to return to the top.\n")
	case statusS3ObjectDeleting:
		w := lipgloss.Width(fmt.Sprintf("%d", m.sum))
		s3keyCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.sum-1)
		spin := m.spinner.View() + " "
		prog := m.progress.View()
		cellsAvail := max(0, m.window.Width-lipgloss.Width(spin+prog+s3keyCount))

		s3keyName := currentNameStyle.Render(m.targetS3Keys[0].String())
		info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Deleting " + s3keyName)
		cellsRemaining := max(0, m.window.Width-lipgloss.Width(spin+info+prog+s3keyCount))
		gap := strings.Repeat(" ", cellsRemaining)
		return spin + info + gap + prog + s3keyCount
	case statusNone, statusS3ObjectFetching:
		return fmt.Sprintf(
			"fetching the list of the S3 objects (profile=%s, bucket=%s)\n",
			m.awsProfile.String(),
			m.bucket.String())
	case statusS3ObjectFetched:
		return m.s3ObjectListString()
	default:
		return m.s3ObjectListString()
	}
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

// getTargetS3Keys returns the list of the S3 bucket objects that the user wants to download
func (m *s3hubListS3ObjectModel) getTargetS3Keys() []model.S3Key {
	targetS3Keys := make([]model.S3Key, 0, len(m.toggles))
	for i, t := range m.toggles {
		if t.Enabled {
			targetS3Keys = append(targetS3Keys, m.s3Keys[i])
		}
	}
	return targetS3Keys
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
	s += ui.Subtle("<space>: choose s3 object to download\n")
	s += ui.Subtle("d: download s3 objects | D: delete s3 objects\n\n")
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
