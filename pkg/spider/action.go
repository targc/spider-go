package spider

type WorkflowAction struct {
	ID         string
	Key        string
	WorkflowID string
	ActionID   string
	Map        map[string]Mapper
}
