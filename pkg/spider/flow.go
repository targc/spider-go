package spider

type FlowStatus string

var (
	FlowStatusDraft  FlowStatus = "draft"
	FlowStatusActive FlowStatus = "active"
)

type Flow struct {
	ID       string            `json:"id"`
	TenantID string            `json:"tenant_id"`
	Name     string            `json:"name"`
	Meta     map[string]string `json:"meta,omitempty"`
	Status   FlowStatus        `json:"status"`
	Version  uint64            `json:"version"`
}
