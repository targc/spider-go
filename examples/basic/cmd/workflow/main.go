package main

import (
	"context"
	"github.com/targc/spider-go/pkg/spider"
	"os"
	"os/signal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	workflow, err := spider.InitDefaultWorkflow(ctx)

	if err != nil {
		panic(err)
	}

	storage := workflow.Storage()

	workflowID := "wa"

	_, err = storage.AddAction(
		ctx,
		workflowID,
		"a1",
		"test-action-a",
		map[string]string{
			"test": "a",
		},
		nil,
	)

	if err != nil {
		panic(err)
	}

	_, err = storage.AddAction(
		ctx,
		workflowID,
		"a2",
		"test-action-b",
		map[string]string{
			"test": "b",
		},
		map[string]spider.Mapper{
			"value": {
				Mode:  spider.MapperModeExpression,
				Value: "a1.output.value + '_updatedx1'",
			},
		},
	)

	if err != nil {
		panic(err)
	}

	_, err = storage.AddAction(
		ctx,
		workflowID,
		"a3",
		"test-action-b",
		map[string]string{
			"test": "c",
		},
		map[string]spider.Mapper{
			"value": {
				Mode:  spider.MapperModeExpression,
				Value: "a2.output.value + '_updatedx2'",
			},
		},
	)

	if err != nil {
		panic(err)
	}

	err = storage.AddDep(ctx, workflowID, "a1", "triggered", "a2")

	if err != nil {
		panic(err)
	}

	err = storage.AddDep(ctx, workflowID, "a2", "success", "a3")

	if err != nil {
		panic(err)
	}

	go workflow.Run(ctx)

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = workflow.Close(ctx)
}
