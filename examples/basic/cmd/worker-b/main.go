package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// ================= WORKER [B] =================

	workerB, err := spider.InitNATSAdapterWorker(ctx, "edge.function_name")

	if err != nil {
		panic(err)
	}

	go workerB.Run(ctx, func(c spider.InputMessageContext, m spider.InputMessage) (*spider.NATSHandlerOutput, error) {

		slog.Info("[process] received input")

		return &spider.NATSHandlerOutput{
			MetaOutput: "success",
			Values: map[string]string{
				"value": "1",
			},
		}, nil
	})

	// ==============================================

	nctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-nctx.Done()

	cancel()
	_ = workerB.Close(ctx)
}
