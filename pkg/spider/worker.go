package spider

import (
	"context"
)

type NATSHandlerOutput struct {
	MetaOutput string
	Values     Values
}

type OutputMessageExternal struct {
	WorkflowNodeID string
	MetaOutput     string
	Values         Values
}

type NATSAdapterWorker struct {
	messenger WorkerMessengerAdapter
	nodeID    string
}

func InitNATSAdapterWorker(
	ctx context.Context,
	nodeID string,
) (*NATSAdapterWorker, error) {
	messenger, err := InitNATSWorkerMessengerAdapter(ctx, nodeID)

	if err != nil {
		return nil, err
	}

	return &NATSAdapterWorker{messenger, nodeID}, nil
}

func (w *NATSAdapterWorker) Run(ctx context.Context, h func(c InputMessageContext, m InputMessage) (*NATSHandlerOutput, error)) error {

	err := w.messenger.ListenInputMessages(
		ctx,
		func(c InputMessageContext, m InputMessage) error {

			output, err := h(c, m)

			if err != nil {
				return err
			}

			err = w.messenger.SendOutputMessage(c.Context, OutputMessage{
				WorkflowNodeID: m.WorkflowNodeID,
				NodeID:         m.NodeID,
				MetaOutput:     output.MetaOutput,
				Values:         output.Values,
			})

			if err != nil {
				return err
			}

			return nil
		},
	)

	return err
}

func (w *NATSAdapterWorker) SendOutputMessage(ctx context.Context, m OutputMessageExternal) error {
	err := w.messenger.SendOutputMessage(ctx, OutputMessage{
		WorkflowNodeID: m.WorkflowNodeID,
		NodeID:         w.nodeID,
		MetaOutput:     m.MetaOutput,
		Values:         m.Values,
	})

	if err != nil {
		return err
	}

	return nil
}

func (w *NATSAdapterWorker) Close(ctx context.Context) error {

	err := w.messenger.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}
