package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// ================= WORKER [A] =================

	workerA, err := spider.InitDefaultWorker(ctx, "test-action-a")

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)

			output := map[string]interface{}{
				"value": "hello",
			}

			outputb, err := json.Marshal(output)

			if err != nil {
				continue
			}

			workerA.SendOutputMessage(ctx, spider.OutputMessage{
				WorkflowActionID: "aaaaa",
				MetaOutput:       "triggered",
				Values:           string(outputb),
			})

			slog.Info("[process] sent")
		}
	}()

	// ==============================================

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = workerA.Close(ctx)
}
