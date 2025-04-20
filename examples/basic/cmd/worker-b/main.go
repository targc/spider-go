package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// ================= WORKER [B] =================

	workerB, err := spider.InitDefaultAdapterWorker(ctx, "test-action-b")

	if err != nil {
		panic(err)
	}

	go workerB.Run(ctx, func(c spider.InputMessageContext, m spider.InputMessage) (*spider.RunOutput, error) {

		slog.Info("[process] received input")

		output := map[string]interface{}{
			"value": "1",
		}

		outputb, err := json.Marshal(output)

		if err != nil {
			return nil, err
		}

		return &spider.RunOutput{
			MetaOutput: "success",
			Values:     outputb,
		}, nil
	})

	// ==============================================

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = workerB.Close(ctx)
}
