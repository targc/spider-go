package spider

type Flow struct {
	ID       string            `json:"id"`
	Version  uint64            `json:"version"`
	Name     string            `json:"name"`
	TenantID string            `json:"tenant_id"`
	Meta     map[string]string `json:"meta,omitempty"`
}
