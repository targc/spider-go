package spider

type Flow struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	TenantID string            `json:"tenant_id"`
	Meta     map[string]string `json:"meta,omitempty"`
}
