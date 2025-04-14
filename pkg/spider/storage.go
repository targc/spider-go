package spider

import "context"

type WorkflowStorageAdapter interface {
	QueryWorkflowNode(ctx context.Context, workflowNodeID string) (*WorkflowNode, error)
	QueryWorkflowNodeDependencies(ctx context.Context, parentWorkflowNodeID, metaOutput string) ([]WorkflowNode, error)
	Close(ctx context.Context) error
}
