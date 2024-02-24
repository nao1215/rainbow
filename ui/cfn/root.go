// Package cfn is the text-based user interface for cfn command.
package cfn

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/ui"
)

// cfnRootModel is the top-level model for the application.
type cfnRootModel struct {
	// err is the error that occurred during the operation.
	err error
	// state is the status of the create bucket operation.
	status status
	// region is the AWS region that the user wants to create the S3 bucket.
	region model.Region
	// awsConfig is the AWS configuration.
	awsConfig *model.AWSConfig
	// awsProfile is the AWS profile.
	awsProfile model.AWSProfile
	// quitting is true when the user has quit the application.
	quitting bool
}

// RunCfnUI start cfn command interactive UI.
func RunCfnUI() error {
	ctx := context.Background()
	profile := model.NewAWSProfile("")
	cfg, err := model.NewAWSConfig(ctx, profile, "")
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(newCFnRootModel(profile, cfg)).Run()
	return err
}

// newCFnRootModel creates a new cfnRootModel.
func newCFnRootModel(profile model.AWSProfile, cfg *model.AWSConfig) *cfnRootModel {
	return &cfnRootModel{
		status:     statusRegionSelecting,
		region:     cfg.Region(),
		awsConfig:  cfg,
		awsProfile: profile,
	}
}

// Init initializes the model.
func (m *cfnRootModel) Init() tea.Cmd {
	return nil
}

// Update is the main update function.
func (m *cfnRootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.region = m.region.Prev()
		case "l", "right":
			m.region = m.region.Next()
		case "enter":
			model, err := newCFnListStackModel(m.region)
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			return model, fetchStacksCmd(model.ctx, model.app, model.region)
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case ui.ErrMsg:
		m.err = msg
		return m, nil
	}
	return m, nil
}

// View renders the model.
func (m cfnRootModel) View() string {
	if m.err != nil {
		return ui.ErrorMessage(m.err)
	}

	if m.quitting {
		return ui.GoodByeMessage()
	}

	return fmt.Sprintf(
		"%s\n\n[ AWS Profile ] %s\n[ ◀︎  %s ▶︎ ] %s\n\n%s\n%s\n",
		"Set the region for the CloudFormation stack you want to display.",
		m.awsProfile.String(),
		ui.Yellow("Region"),
		ui.Green(m.region.String()),
		ui.Subtle("h/l, left/right: select region | <esc>, <Ctrl-C>, q: quit"),
		ui.Subtle("<enter>: list up the CloudFormation stacks"))
}
