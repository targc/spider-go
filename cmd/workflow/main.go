package main

import (
	"context"
	"github.com/targc/spider-go/pkg/spider"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorkflowAction struct {
	Key      string                   `json:"key"`
	ActionID string                   `json:"action_id"`
	Config   map[string]string        `json:"config"`
	Mapper   map[string]spider.Mapper `json:"mapper"`
}

type Peer struct {
	ParentKey  string `json:"parent_key"`
	MetaOutput string `json:"meta_output"`
	ChildKey   string `json:"child_key"`
}

type WorkflowCreatePayload struct {
	Actions []WorkflowAction `json:"actions"`
	Peers   []Peer           `json:"peers"`
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

		// TODO: validate graph & input mapper schema

		for _, action := range payload.Actions {
			_, err = storage.AddAction(
				ctx,
				workflowID,
				action.Key,
				action.ActionID,
				action.Config,
				action.Mapper,
			)

			if err != nil {
				return err
			}
		}

		for _, peer := range payload.Peers {
			err = storage.AddDep(
				ctx,
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
