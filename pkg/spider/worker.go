package spider

import (
	"context"
)

type Worker struct {
	messenger WorkerMessengerAdapter
	actionID  string
}

func InitDefaultWorker(
	ctx context.Context,
	actionID string,
) (*Worker, error) {
	messenger, err := InitNATSWorkerMessengerAdapter(ctx, actionID)

	if err != nil {
		return nil, err
	}

	return &Worker{messenger, actionID}, nil
}

func (w *Worker) Run(ctx context.Context, h func(c InputMessageContext, m InputMessage) error) error {

	err := w.messenger.ListenInputMessages(
		ctx,
		func(c InputMessageContext, m InputMessage) error {

			c.SendOutput = func(metaOutput string, values string) error {
				err := w.SendOutputMessage(c.Context, OutputMessage{
					SessionID:        m.SessionID,
					WorkflowActionID: m.WorkflowActionID,
					MetaOutput:       metaOutput,
					Values:           values,
				})

				if err != nil {
					return err
				}

				return nil
			}

			err := h(c, m)

			if err != nil {
				return err
			}

			return nil
		},
	)

	return err
}

func (w *Worker) SendOutputMessage(ctx context.Context, m OutputMessage) error {
	err := w.messenger.SendOutputMessage(ctx, OutputMessage{
		SessionID:        m.SessionID,
		WorkflowActionID: m.WorkflowActionID,
		MetaOutput:       m.MetaOutput,
		Values:           m.Values,
	})

	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Close(ctx context.Context) error {

	err := w.messenger.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}
