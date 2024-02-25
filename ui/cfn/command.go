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
	events map[string][]*model.StackEvent
}

func fetchStacksCmd(ctx context.Context, app *di.CFnApp, region model.Region) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		listCFnStackOutput, err := app.CFnStackLister.ListCFnStack(ctx, &usecase.CFnStackListerInput{
			Region: region,
		})
		if err != nil {
			return ui.ErrMsg(err)
		}

		stacks := make([]*model.Stack, 0, len(listCFnStackOutput.Stacks))
		for _, stack := range listCFnStackOutput.Stacks {
			if stack.StackName == nil || stack.StackStatus == model.StackStatusDeleteComplete {
				continue
			}
			stacks = append(stacks, stack)
		}

		events := make(map[string][]*model.StackEvent)
		for _, stack := range stacks {
			describeCFnStackEventsOutput, err := app.CFnStackEventsDescriber.DescribeCFnStackEvents(ctx, &usecase.CFnStackEventsDescriberInput{
				StackName: *stack.StackName,
				Region:    region,
			})
			if err != nil {
				return ui.ErrMsg(err)
			}
			events[*stack.StackName] = describeCFnStackEventsOutput.Events
		}

		return fetchStacks{
			stacks: stacks,
			events: events,
		}
	})
}
