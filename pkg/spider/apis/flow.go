package apis

import (
	"github.com/targc/spider-go/pkg/spider/usecase"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateFlow(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	var payload struct {
		Name    string            `json:"name"`
		Meta    map[string]string `json:"meta,omitempty"`
		Actions []WorkflowAction  `json:"actions"`
		Peers   []Peer            `json:"peers"`
	}

	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	if payload.Name == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "name is required",
		})
	}

	// Convert to usecase types
	actions := make([]usecase.WorkflowActionInput, len(payload.Actions))
	for i, action := range payload.Actions {
		actions[i] = usecase.WorkflowActionInput{
			Key:      action.Key,
			ActionID: action.ActionID,
			Config:   action.Config,
			Mapper:   action.Mapper,
			Meta:     action.Meta,
		}
	}

	peers := make([]usecase.PeerInput, len(payload.Peers))
	for i, peer := range payload.Peers {
		peers[i] = usecase.PeerInput{
			ParentKey:  peer.ParentKey,
			MetaOutput: peer.MetaOutput,
			ChildKey:   peer.ChildKey,
		}
	}

	req := &usecase.CreateFlowRequest{
		TenantID: tenantID,
		Name:     payload.Name,
		Meta:     payload.Meta,
		Actions:  actions,
		Peers:    peers,
	}

	result, err := h.usecase.CreateFlow(req)
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"error": "Failed to create flow",
		})
	}

	return c.JSON(result)
}

func (h *Handler) ListFlows(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	pageSize := c.QueryInt("page_size", 20)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.usecase.ListFlows(tenantID, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"error": "Failed to list flows",
		})
	}

	return c.JSON(result)
}

func (h *Handler) GetFlow(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	flowID := c.Params("id")
	if flowID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "flow id is required",
		})
	}

	result, err := h.usecase.GetFlow(tenantID, flowID)
	if err != nil {
		return c.Status(404).JSON(map[string]string{
			"error": "Flow not found",
		})
	}

	return c.JSON(result)
}

func (h *Handler) DeleteFlow(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	flowID := c.Params("flow_id")
	if flowID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "flow_id is required",
		})
	}

	err := h.usecase.DeleteFlow(tenantID, flowID)
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"error": "Failed to delete flow",
		})
	}

	return c.Status(204).Send(nil)
}