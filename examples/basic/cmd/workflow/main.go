package main

import (
	"context"
	"errors"
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
			Key:        "test_workflow_action_a",
			WorkflowID: "z",
			ActionID:   "test-action-a",
		}, nil
	case "bbbbb":
		return &spider.WorkflowAction{
			ID:         workflowActionID,
			Key:        "test_workflow_action_b",
			WorkflowID: "z",
			ActionID:   "test-action-b",
		}, nil
	}

	return nil, errors.New("not found workflow action")
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowActionDependencies(ctx context.Context, parentWorkflowActionID, metaOutput string) ([]spider.WorkflowAction, error) {

	switch parentWorkflowActionID {
	case "aaaaa":
		return []spider.WorkflowAction{
			{
				ID:         "bbbbb",
				Key:        "test_workflow_action_b",
				WorkflowID: "test-workflow-a",
				ActionID:   "test-action-b",
			},
		}, nil
	}

	return nil, errors.New("not found deps")
}

func (w *MockWorkflowStorageAdapter) QueryWorkflowActionMapper(ctx context.Context, parentWorkflowActionID, metaOutput, workflowActionID string) (map[string]spider.Mapper, error) {

	switch parentWorkflowActionID {
	case "aaaaa":
		switch metaOutput {
		case "triggered":
			switch workflowActionID {
			case "bbbbb":
				return map[string]spider.Mapper{
					"mapped_test_1": {
						Mode:  spider.MapperModeFixed,
						Value: "hello_fixed",
					},
					"mapped_test_2": {
						Mode:  spider.MapperModeKey,
						Value: "test_workflow_action_a.output.value",
					},
					"mapped_test_3": {
						Mode:  spider.MapperModeExpression,
						Value: "test_workflow_action_a.output.value + '_suffix'",
					},
				}, nil
			}
		}
	}

	return nil, errors.New("not found mapper")
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
