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

	workerA, err := spider.InitDefaultAdapterWorker(ctx, "test-action-a")

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)

			workerA.SendOutputMessage(ctx, spider.OutputMessageExternal{
				WorkflowActionID: "test-workflow-action-a",
				MetaOutput:       "triggered",
				Values: map[string]interface{}{
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
