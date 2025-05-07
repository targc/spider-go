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
	QueryWorkflowAction(ctx context.Context, workflowID, key string) (*WorkflowAction, error)
	QueryWorkflowActionDependencies(ctx context.Context, workflowID, key, metaOutput string) ([]WorkflowAction, error)
	QueryWorkflowActionMapper(ctx context.Context, workflowID, key, metaOutput, depKey string) (map[string]Mapper, error)
	AddAction(ctx context.Context, workflowID, key, actionID string, config map[string]string) (*WorkflowAction, error)
	AddDep(ctx context.Context, workflowID, key, metaOutput, key2 string, m map[string]Mapper) error
	TryAddSessionContext(ctx context.Context, workflowID, sessionID, key string, value map[string]interface{}) (map[string]map[string]interface{}, error)
	Close(ctx context.Context) error
}
