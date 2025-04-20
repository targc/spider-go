package spider

import "context"

type WorkflowStorageAdapter interface {
	QueryWorkflowAction(ctx context.Context, workflowActionID string) (*WorkflowAction, error)
	QueryWorkflowActionDependencies(ctx context.Context, parentWorkflowActionID, metaOutput string) ([]WorkflowAction, error)
	Close(ctx context.Context) error
}
