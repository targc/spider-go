package spider

type WorkflowAction struct {
	ID         string
	Key        string
	TenantID   string
	WorkflowID string
	ActionID   string
	Config     map[string]string
	Map        map[string]Mapper
	Disabled   bool
}
