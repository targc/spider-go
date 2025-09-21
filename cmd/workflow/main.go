package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/targc/spider-go/pkg/spider"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorkflowAction struct {
	Key      string                   `json:"key"`
	ActionID string                   `json:"action_id"`
	Config   map[string]string        `json:"config"`
	Mapper   map[string]spider.Mapper `json:"mapper"`
	Meta     map[string]string        `json:"meta,omitempty"`
}

type Peer struct {
	ParentKey  string `json:"parent_key"`
	MetaOutput string `json:"meta_output"`
	ChildKey   string `json:"child_key"`
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	worflow, err := spider.InitDefaultWorkflow(ctx)

	if err != nil {
		panic(err)
	}

	storage := worflow.Storage()

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return nil
	})

	app.Post("/tenants/:tenant_id/workflows", func(c *fiber.Ctx) error {

		tenantID := c.Params("tenant_id")
		if tenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id is required",
			})
		}

		var payload struct {
			Name    string           `json:"name"`
			Actions []WorkflowAction `json:"actions"`
			Peers   []Peer           `json:"peers"`
		}

		err = c.BodyParser(&payload)

		if err != nil {
			return err
		}

		if payload.Name == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "name is required",
			})
		}

		id, err := uuid.NewV7()
		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to generate workflow ID",
			})
		}

		workflow, err := storage.CreateWorkflow(ctx, &spider.CreateWorkflowRequest{
			ID:       id.String(),
			TenantID: tenantID,
			Name:     payload.Name,
		})

		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to create workflow",
			})
		}

		workflowID := workflow.ID

		// TODO: validate graph & input mapper schema

		for _, action := range payload.Actions {
			_, err = storage.AddAction(ctx, &spider.AddActionRequest{
				TenantID:   tenantID,
				WorkflowID: workflowID,
				Key:        action.Key,
				ActionID:   action.ActionID,
				Config:     action.Config,
				Map:        action.Mapper,
				Meta:       action.Meta,
			})

			if err != nil {
				return err
			}
		}

		for _, peer := range payload.Peers {
			err = storage.AddDep(
				ctx,
				tenantID,
				workflowID,
				peer.ParentKey,
				peer.MetaOutput,
				peer.ChildKey,
			)

			if err != nil {
				return err
			}
		}

		return c.JSON(map[string]interface{}{
			"workflow_id":   workflowID,
			"workflow_name": workflow.Name,
		})
	})

	app.Get("/tenants/:tenant_id/workflows", func(c *fiber.Ctx) error {
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

		result, err := storage.ListWorkflows(ctx, tenantID, page, pageSize)
		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to list workflows",
			})
		}

		return c.JSON(result)
	})

	app.Get("/tenants/:tenant_id/workflows/:id", func(c *fiber.Ctx) error {
		tenantID := c.Params("tenant_id")
		if tenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id is required",
			})
		}

		workflowID := c.Params("id")
		if workflowID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "workflow id is required",
			})
		}

		workflow, err := storage.GetWorkflow(ctx, tenantID, workflowID)
		if err != nil {
			return c.Status(404).JSON(map[string]string{
				"error": "Workflow not found",
			})
		}

		actions, err := storage.GetWorkflowActions(ctx, tenantID, workflowID)
		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to get workflow actions",
			})
		}

		return c.JSON(map[string]interface{}{
			"workflow_id":   workflowID,
			"workflow_name": workflow.Name,
			"tenant_id":     tenantID,
			"actions":       actions,
		})
	})

	app.Post("/tenants/:tenant_id/workflows/:workflow_id/actions/:key/disable", func(c *fiber.Ctx) error {

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

		err = storage.DisableWorkflowAction(ctx, tenantID, workflowID, key)

		if err != nil {
			return err
		}

		slog.Info("[process] disabled")

		return nil
	})

	app.Put("/tenants/:tenant_id/workflows/:workflow_id/actions/:key", func(c *fiber.Ctx) error {

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

		err = c.BodyParser(&payload)

		if err != nil {
			return err
		}

		action, err := storage.UpdateAction(ctx, &spider.UpdateActionRequest{
			TenantID:   tenantID,
			WorkflowID: workflowID,
			Key:        key,
			Config:     payload.Config,
			Map:        payload.Mapper,
			Meta:       payload.Meta,
		})

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
	})

	app.Delete("/tenants/:tenant_id/workflows/:workflow_id", func(c *fiber.Ctx) error {

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

		err = storage.DeleteWorkflow(ctx, tenantID, workflowID)

		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to delete workflow",
			})
		}

		return c.Status(204).Send(nil)
	})

	go worflow.Run(ctx)

	go func() {
		err := app.Listen("0.0.0.0:8080")

		if err != nil {
			panic(err)
		}
	}()

	nctx, ncancel := signal.NotifyContext(ctx, os.Interrupt)
	defer ncancel()

	<-nctx.Done()

	cancel()
	_ = worflow.Close(ctx)
}
