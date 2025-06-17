package spider

import (
	"context"
	"log/slog"
)

type Worker struct {
	messenger WorkerMessengerAdapter
	storage   WorkerStorageAdapter
	actionID  string
}

func InitDefaultWorker(
	ctx context.Context,
	actionID string,
) (*Worker, error) {
	messenger, err := InitNATSWorkerMessengerAdapter(ctx, actionID, InitNATSWorkerMessengerAdapterOpt{
		BetaAutoSetupNATS: true,
	})

	if err != nil {
		return nil, err
	}

	storage, err := InitMongodDBWorkerStorageAdapter(ctx)

	if err != nil {
		return nil, err
	}

	return &Worker{
		messenger,
		storage,
		actionID,
	}, nil
}

func (w *Worker) Run(ctx context.Context, h func(c InputMessageContext, m InputMessage) error) error {

	err := w.messenger.ListenInputMessages(
		ctx,
		func(c InputMessageContext, m InputMessage) error {

			c.SendOutput = func(metaOutput string, values string) error {
				err := w.messenger.SendOutputMessage(c.Context, m.ToOutputMessage(metaOutput, values))

				if err != nil {
					return err
				}

				return nil
			}

			err := h(c, m)

			if err != nil {
				slog.Error("failed to process handler", slog.String("error", err.Error()))
				return err
			}

			return nil
		},
	)

	return err
}

func (w *Worker) SendTriggerMessage(ctx context.Context, m TriggerMessage) error {

	m.ActionID = w.actionID

	err := w.messenger.SendTriggerMessage(ctx, m)

	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) GetAllConfigs(ctx context.Context) ([]WorkerConfig, error) {

	confs, err := w.storage.GetAllConfigs(ctx, w.actionID)

	if err != nil {
		return nil, err
	}

	return confs, nil
}

func (w *Worker) Close(ctx context.Context) error {

	err := w.messenger.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}
