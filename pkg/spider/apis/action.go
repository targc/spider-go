package apis

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/targc/spider-go/pkg/spider"
)

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
