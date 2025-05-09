package spider

import (
	"context"
	"os"
	"os/signal"
)

func LazyBootstrapWorker(actionID string, h func(c InputMessageContext, m InputMessage) error) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerB, err := InitDefaultWorker(ctx, actionID)

	if err != nil {
		return err
	}

	go workerB.Run(ctx, h)

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()

	_ = workerB.Close(ctx)

	return nil
}
