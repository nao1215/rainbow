package cfn

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nao1215/rainbow/app/di"
	"github.com/nao1215/rainbow/app/domain/model"
	"github.com/nao1215/rainbow/app/usecase"
	"github.com/nao1215/rainbow/ui"
)

// fetchStacks is the message that is sent when the user wants to fetch the list of the CloudFormation stacks.
type fetchStacks struct {
	stacks []*model.Stack
}

func fetchStacksCmd(ctx context.Context, app *di.CFnApp, region model.Region) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		output, err := app.CFnStackLister.ListCFnStack(ctx, &usecase.CFnStackListerInput{
			Region: region,
		})
		if err != nil {
			return ui.ErrMsg(err)
		}
		return fetchStacks{
			stacks: output.Stacks,
		}
	})
}
