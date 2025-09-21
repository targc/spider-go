package apis

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/targc/spider-go/pkg/spider"
)

// DisableAction godoc
// @Summary Disable a workflow action
// @Description Disable a specific workflow action
// @Tags actions
// @Param tenant_id path string true "Tenant ID"
// @Param workflow_id path string true "Workflow ID"
// @Param key path string true "Action Key"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/workflows/{workflow_id}/actions/{key}/disable [post]
func (h *Handler) DisableAction(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	workflowID := c.Params("workflow_id")
	if workflowID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "workflow_id is required",
		})
	}

	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "key is required",
		})
	}

	err := h.usecase.DisableAction(c.Context(), tenantID, workflowID, key)
	if err != nil {
		return err
	}

	slog.Info("[process] disabled")

	return nil
}

// UpdateAction godoc
// @Summary Update a workflow action
// @Description Update configuration of a workflow action
// @Tags actions
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param workflow_id path string true "Workflow ID"
// @Param key path string true "Action Key"
// @Param payload body object true "Action update payload"
// @Success 200 {object} spider.WorkflowAction
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/workflows/{workflow_id}/actions/{key} [put]
func (h *Handler) UpdateAction(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	if tenantID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "tenant_id is required",
		})
	}

	workflowID := c.Params("workflow_id")
	if workflowID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "workflow_id is required",
		})
	}

	key := c.Params("key")
	if key == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "key is required",
		})
	}

	var payload struct {
		Config map[string]string        `json:"config"`
		Mapper map[string]spider.Mapper `json:"mapper"`
		Meta   map[string]string        `json:"meta,omitempty"`
	}

	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	req := &spider.UpdateActionRequest{
		TenantID:   tenantID,
		WorkflowID: workflowID,
		Key:        key,
		Config:     payload.Config,
		Map:        payload.Mapper,
		Meta:       payload.Meta,
	}

	action, err := h.usecase.UpdateAction(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(map[string]string{
			"error": "Failed to update action",
		})
	}

	if action == nil {
		return c.Status(404).JSON(map[string]string{
			"error": "Action not found",
		})
	}

	return c.JSON(action)
}
