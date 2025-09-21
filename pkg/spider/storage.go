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

type WorkflowInfo struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
}

type WorkflowListResponse struct {
	Workflows []WorkflowInfo `json:"workflows"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
}

type WorkflowStorageAdapter interface {
	QueryWorkflowAction(ctx context.Context, tenantID, workflowID, key string) (*WorkflowAction, error)
	QueryWorkflowActionDependencies(ctx context.Context, tenantID, workflowID, key, metaOutput string) ([]WorkflowAction, error)
	AddAction(ctx context.Context, tenantID, workflowID, key, actionID string, conf map[string]string, m map[string]Mapper, meta map[string]string) (*WorkflowAction, error)
	AddDep(ctx context.Context, tenantID, workflowID, key, metaOutput, key2 string) error
	GetSessionContext(ctx context.Context, workflowID, sessionID, taskID string) (map[string]map[string]interface{}, error)
	CreateSessionContext(ctx context.Context, workflowID, sessionID, taskID string, value map[string]map[string]interface{}) error
	DeleteSessionContext(ctx context.Context, workflowID, sessionID, taskID string) error
	DisableWorkflowAction(ctx context.Context, tenantID, workflowID, key string) error
	ListWorkflows(ctx context.Context, tenantID string, page, pageSize int) (*WorkflowListResponse, error)
	GetWorkflowActions(ctx context.Context, tenantID, workflowID string) ([]WorkflowAction, error)
	Close(ctx context.Context) error
}

type WorkerConfig struct {
	WorkflowActionID string            `json:"workflow_action_id"`
	TenantID         string            `json:"tenant_id"`
	WorkflowID       string            `json:"workflow_id"`
	Key              string            `json:"key"`
	Config           map[string]string `json:"config"`
	Meta             map[string]string `json:"meta,omitempty"`
}

type WorkerStorageAdapter interface {
	GetAllConfigs(ctx context.Context, actionID string) ([]WorkerConfig, error)
	Close(ctx context.Context) error
}
