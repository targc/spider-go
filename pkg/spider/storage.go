package spider

import "context"

type MapperMode string

var (
	MapperModeFixed      MapperMode = "fixed"
	MapperModeKey        MapperMode = "key"
	MapperModeExpression MapperMode = "expression"
)

type Mapper struct {
	Mode  MapperMode
	Value string
}

type WorkflowStorageAdapter interface {
	QueryWorkflowAction(ctx context.Context, workflowActionID string) (*WorkflowAction, error)
	QueryWorkflowActionDependencies(ctx context.Context, parentWorkflowActionID, metaOutput string) ([]WorkflowAction, error)
	QueryWorkflowActionMapper(ctx context.Context, parentWorkflowActionID, metaOutput, workflowActionID string) (map[string]Mapper, error)
	TryAddSessionContext(ctx context.Context, sessionID, contextKey string, contextValue map[string]interface{}) (map[string]interface{}, error)
	Close(ctx context.Context) error
}
