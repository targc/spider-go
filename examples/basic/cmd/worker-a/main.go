package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// ================= WORKER [A] =================

	workerA, err := spider.InitNATSAdapterWorker(ctx, "test-node-a")

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)

			workerA.SendOutputMessage(ctx, spider.OutputMessageExternal{
				WorkflowNodeID: "test-workflow-node-a",
				MetaOutput:     "triggered",
				Values: map[string]string{
					"value": "hello",
				},
			})

			slog.Info("[process] sent")
		}
	}()

	// ==============================================

	nctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-nctx.Done()

	cancel()
	_ = workerA.Close(ctx)
}
