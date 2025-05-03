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

	workerB, err := spider.InitDefaultWorker(ctx, "test-action-b")

	if err != nil {
		panic(err)
	}

	go workerB.Run(ctx, func(c spider.InputMessageContext, m spider.InputMessage) error {

		slog.Info("[process] received input", slog.Any("message", m))

		var input struct {
			Value string `json:"value"`
		}

		err := json.Unmarshal([]byte(m.Values), &input)

		if err != nil {
			return err
		}

		output := map[string]interface{}{
			"value": input.Value,
		}

		outputb, err := json.Marshal(output)

		if err != nil {
			return err
		}

		err = c.SendOutput("success", string(outputb))

		if err != nil {
			return err
		}

		return nil
	})

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = workerB.Close(ctx)
}
