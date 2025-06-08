package spider

import "context"

type MapperMode string

var (
	MapperModeFixed      MapperMode = "fixed"
	MapperModeKey        MapperMode = "key"
	MapperModeExpression MapperMode = "expression"
)

type Mapper struct {
	Mode  MapperMode `json:"mode"`
	Value string     `json:"value"`
}

type WorkflowStorageAdapter interface {
	QueryWorkflowAction(ctx context.Context, workflowID, key string) (*WorkflowAction, error)
	QueryWorkflowActionDependencies(ctx context.Context, workflowID, key, metaOutput string) ([]WorkflowAction, error)
	AddAction(ctx context.Context, workflowID, key, actionID string, conf map[string]string, m map[string]Mapper) (*WorkflowAction, error)
	AddDep(ctx context.Context, workflowID, key, metaOutput, key2 string) error
	GetSessionContext(ctx context.Context, workflowID, sessionID, taskID string) (map[string]map[string]interface{}, error)
	CreateSessionContext(ctx context.Context, workflowID, sessionID, taskID string, value map[string]map[string]interface{}) error
	DeleteSessionContext(ctx context.Context, workflowID, sessionID, taskID string) error
	Close(ctx context.Context) error
}
