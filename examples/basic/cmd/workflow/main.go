package main

import (
	"context"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
)

type MockWorkflowStorageAdapter struct {
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowAction(ctx context.Context, workflowActionID string) (*spider.WorkflowAction, error) {

	switch workflowActionID {
	case "aaaaa":
		return &spider.WorkflowAction{
			ID:         workflowActionID,
			Key:        "test-workflow-action-a",
			WorkflowID: "z",
			ActionID:   "test-action-a",
		}, nil
	}

	return nil, nil
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowActionDependencies(ctx context.Context, parentWorkflowActionID, metaOutput string) ([]spider.WorkflowAction, error) {

	switch parentWorkflowActionID {
	case "aaaaa":
		return []spider.WorkflowAction{
			{
				ID:         "bbbbb",
				Key:        "test-workflow-action-b",
				WorkflowID: "test-workflow-a",
				ActionID:   "test-action-b",
			},
		}, nil
	}

	return nil, nil
}

func (w *MockWorkflowStorageAdapter) Close(ctx context.Context) error {
	return nil
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// ================= WORKFLOW =================

	mockStorage := MockWorkflowStorageAdapter{}

	worflow, err := spider.InitDefaultWorkflow(ctx, &mockStorage)

	if err != nil {
		panic(err)
	}

	go worflow.Run(ctx)

	// ==============================================

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = worflow.Close(ctx)
}
