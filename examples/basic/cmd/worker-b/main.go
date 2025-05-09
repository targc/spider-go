package main

import (
	"encoding/json"
	"log/slog"
	"spider-go/pkg/spider"
)

const actionID = "test-action-b"

func main() {
	err := spider.LazyBootstrapWorker(actionID, func(c spider.InputMessageContext, m spider.InputMessage) error {

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

	if err != nil {
		panic(err)
	}
}
