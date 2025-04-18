package spider

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Workflow struct {
	messenger WorkflowMessengerAdapter
	storage   WorkflowStorageAdapter
}

func InitWorkflow(
	messenger WorkflowMessengerAdapter,
	storage WorkflowStorageAdapter,
) *Workflow {
	return &Workflow{
		messenger,
		storage,
	}
}

func InitDefaultWorkflow(
	ctx context.Context,
	storage WorkflowStorageAdapter,
) (*Workflow, error) {
	messenger, err := InitNATSWorkflowMessengerAdapter(ctx)

	if err != nil {
		return nil, err
	}

	return &Workflow{
		messenger,
		storage,
	}, nil
}

func (w *Workflow) Run(ctx context.Context) error {

	err := w.messenger.ListenOutputMessages(ctx, func(c OutputMessageContext, m OutputMessage) error {

		workflowNode, err := w.storage.QueryWorkflowNode(c.Context, m.WorkflowNodeID)

		if err != nil {
			return err
		}

		_ = workflowNode // TODO:

		deps, err := w.storage.QueryWorkflowNodeDependencies(c.Context, m.WorkflowNodeID, m.MetaOutput)

		if err != nil {
			return err
		}

		eg := errgroup.Group{}

		eg.SetLimit(10)

		for _, dep := range deps {
			eg.Go(func() error {
				err = w.messenger.SendInputMessage(ctx, InputMessage{
					WorkflowNodeID: dep.ID,
					NodeID:         dep.NodeID,
					Values:         m.Values, // TODO: transformer, value mapper
				})

				if err != nil {
					return err
				}

				return nil
			})
		}

		_ = eg.Wait()

		return nil
	})

	return err
}

func (w *Workflow) Close(ctx context.Context) error {

	err := w.messenger.Close(ctx)

	if err != nil {
		return err
	}

	err = w.storage.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}
