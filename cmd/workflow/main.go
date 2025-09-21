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

type WorkflowCreatePayload struct {
	TenantID string           `json:"tenant_id"`
	Actions  []WorkflowAction `json:"actions"`
	Peers    []Peer           `json:"peers"`
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

	app.Post("/workflows", func(c *fiber.Ctx) error {

		var payload WorkflowCreatePayload

		err = c.BodyParser(&payload)

		if err != nil {
			return err
		}

		workflowUUID, err := uuid.NewV7()

		if err != nil {
			return err
		}

		workflowID := workflowUUID.String()

		if payload.TenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id is required",
			})
		}

		tenantID := payload.TenantID

		// TODO: validate graph & input mapper schema

		for _, action := range payload.Actions {
			_, err = storage.AddAction(
				ctx,
				tenantID,
				workflowID,
				action.Key,
				action.ActionID,
				action.Config,
				action.Mapper,
				action.Meta,
			)

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
			"workflow_id": workflowID,
		})
	})

	app.Get("/workflows", func(c *fiber.Ctx) error {
		tenantID := c.Query("tenant_id")
		if tenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id query parameter is required",
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

	app.Get("/workflows/:id", func(c *fiber.Ctx) error {
		workflowID := c.Params("id")
		if workflowID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "workflow id is required",
			})
		}

		tenantID := c.Query("tenant_id")
		if tenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id query parameter is required",
			})
		}

		actions, err := storage.GetWorkflowActions(ctx, tenantID, workflowID)
		if err != nil {
			return c.Status(500).JSON(map[string]string{
				"error": "Failed to get workflow",
			})
		}

		if len(actions) == 0 {
			return c.Status(404).JSON(map[string]string{
				"error": "Workflow not found",
			})
		}

		return c.JSON(map[string]interface{}{
			"workflow_id": workflowID,
			"tenant_id":   tenantID,
			"actions":     actions,
		})
	})

	app.Post("/workflow-action-disable", func(c *fiber.Ctx) error {

		var payload struct {
			TenantID   string `json:"tenant_id"`
			WorkflowID string `json:"workflow_id"`
			Key        string `json:"key"`
		}

		err = c.BodyParser(&payload)

		if err != nil {
			return err
		}

		if payload.TenantID == "" {
			return c.Status(400).JSON(map[string]string{
				"error": "tenant_id is required",
			})
		}

		err = storage.DisableWorkflowAction(ctx, payload.TenantID, payload.WorkflowID, payload.Key)

		if err != nil {
			return err
		}

		slog.Info("[process] disabled")

		return nil
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
