package main

import (
	"context"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	worflow, err := spider.InitDefaultWorkflow(ctx)

	if err != nil {
		panic(err)
	}

	storage := worflow.Storage()

	workflowID := "wa"

	_, err = storage.AddAction(ctx, workflowID, "a1", "test-action-a", nil)

	if err != nil {
		panic(err)
	}

	_, err = storage.AddAction(ctx, workflowID, "a2", "test-action-b", map[string]spider.Mapper{
		"value": {
			Mode:  spider.MapperModeExpression,
			Value: "a1.output.value + '_updatedx1'",
		},
	})

	if err != nil {
		panic(err)
	}

	_, err = storage.AddAction(ctx, workflowID, "a3", "test-action-b", map[string]spider.Mapper{
		"value": {
			Mode:  spider.MapperModeExpression,
			Value: "a2.output.value + '_updatedx2'",
		},
	})

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

	go worflow.Run(ctx)

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = worflow.Close(ctx)
}
