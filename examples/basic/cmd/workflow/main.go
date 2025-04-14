package main

import (
	"context"
	"os"
	"os/signal"
	"spider-go/pkg/spider"
)

type MockWorkflowStorageAdapter struct {
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowNode(ctx context.Context, workflowNodeID string) (*spider.WorkflowNode, error) {

	switch workflowNodeID {
	case "test-workflow-node-a":
		return &spider.WorkflowNode{
			ID:         workflowNodeID,
			WorkflowID: "test-workflow-a",
			NodeID:     "test-node-a",
		}, nil
	}

	return nil, nil
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowNodeDependencies(ctx context.Context, parentWorkflowNodeID, metaOutput string) ([]spider.WorkflowNode, error) {

	switch parentWorkflowNodeID {
	case "test-workflow-node-a":
		return []spider.WorkflowNode{
			{
				ID:         "test-workflow-node-b",
				WorkflowID: "test-workflow-a",
				NodeID:     "test-node-b",
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

	nctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-nctx.Done()

	cancel()
	_ = worflow.Close(ctx)
}
