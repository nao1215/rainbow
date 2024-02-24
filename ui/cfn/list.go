package cfn

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/ui"
)

// cfnListStackModel is the model for listing the CloudFormation stacks.
type cfnListStackModel struct {
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
	// app is the CFn application service.
	app *di.CFnApp
	// ctx is the context.
	ctx context.Context
	// stacks is the list of the CloudFormation stacks.
	stacks []*model.Stack
	// status is the status of the operation.
	status status
	// toggle is the currently selected menu item.
	toggles ui.ToggleSets
}

const (
	windowHeight = 10
)

// newCFnListStackModel returns the new cfnListStackModel for listing the CloudFormation stacks.
func newCFnListStackModel(region model.Region) (*cfnListStackModel, error) {
	ctx := context.Background()
	profile := model.NewAWSProfile("")
	cfg, err := model.NewAWSConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	app, err := di.NewCFnApp(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	return &cfnListStackModel{
		awsConfig:  cfg,
		awsProfile: profile,
		region:     region,
		app:        app,
		stacks:     []*model.Stack{},
		status:     statusNone,
		choice:     ui.NewChoice(0, 0),
		ctx:        ctx,
	}, nil
}

// Init initializes the model.
func (m *cfnListStackModel) Init() tea.Cmd {
	return nil // Not called this method
}

// Update updates the model.
func (m *cfnListStackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return newCFnRootModel(m.awsProfile, m.awsConfig), nil
		case "D":
			// TODO: implement delete stack
		case "enter":
			if m.status == statusReturnToTop {
				return newCFnRootModel(m.awsProfile, m.awsConfig), nil
			}
		case " ":
			if m.status == statusStacksListed {
				m.toggles[m.choice.Choice].Toggle()
			}
		}
	case fetchStacks:
		m.status = statusStacksFetched
		m.stacks = make([]*model.Stack, 0, len(msg.stacks))
		for _, stack := range msg.stacks {
			if stack.StackName == nil || stack.StackStatus == model.StackStatusDeleteComplete {
				continue
			}
			m.stacks = append(m.stacks, stack)
		}
		m.choice = ui.NewChoice(0, len(m.stacks)-1)
		m.toggles = ui.NewToggleSets(len(m.stacks))
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
func (m *cfnListStackModel) View() string {
	if m.err != nil {
		m.status = statusQuit
		return ui.ErrorMessage(m.err)
	}

	switch m.status {
	case statusQuit:
		return ui.GoodByeMessage()
	case statusNone, statusStacksFetching:
		return fmt.Sprintf(
			"fetching the list of the CloudForamtion Stack (profile=%s)\n",
			m.awsProfile.String())
	case statusStacksFetched:
		return m.stackListString()
	default:
		return m.stackListString()
	}
}

// stackListString returns the string of the list of the CloudFormation stacks.
func (m *cfnListStackModel) stackListString() string {
	switch len(m.stacks) {
	case 0:
		return m.emptyStacksString()
	default:
		return m.stacksListStrWithCheckBox()
	}
}

// emptyStacksString returns the string of the empty list of the CloudFormation stacks.
func (m *cfnListStackModel) emptyStacksString() string {
	m.status = statusReturnToTop
	return fmt.Sprintf("No CloudFormation Stacks (profile=%s)\n\n%s\n",
		m.awsProfile.String(),
		ui.Subtle("<enter>: return to the top"))
}

// stacksListStrWithCheckBox returns the string of the list of the CloudFormation stacks with checkbox.
func (m *cfnListStackModel) stacksListStrWithCheckBox() string {
	startIndex := 0
	endIndex := len(m.stacks)

	if m.choice.Choice >= windowHeight {
		startIndex = m.choice.Choice - windowHeight + 1
		endIndex = startIndex + windowHeight
		if endIndex > len(m.stacks) {
			startIndex = len(m.stacks) - windowHeight
			endIndex = len(m.stacks)
		}
	} else {
		if len(m.stacks) > windowHeight {
			endIndex = windowHeight
		}
	}

	m.status = statusStacksListed
	s := fmt.Sprintf("CloudForamtion Stacks %d/%d (profile=%s)\n\n", m.choice.Choice+1, len(m.stacks), m.awsProfile.String())
	for i := startIndex; i < endIndex; i++ {
		stack := m.stacks[i]
		s += fmt.Sprintf("%s\n",
			ui.ToggleWidget(
				fmt.Sprintf(
					"  %s (status=%s, updated_at=%s)",
					color.GreenString(*stack.StackName),
					stack.StackStatus.StringWithColor(),
					stack.LastUpdatedTime.Format("2006-01-02 15:04:05")),
				m.choice.Choice == i, m.toggles[i].Enabled))
	}
	s += ui.Subtle("\n<esc>: return to the top | <Ctrl-C>: quit | up/down: select\n")
	return s
}
