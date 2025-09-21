package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/targc/spider-go/pkg/spider"
	"github.com/targc/spider-go/pkg/spider/apis"
	"github.com/targc/spider-go/pkg/spider/usecase"

	"github.com/gofiber/fiber/v2"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	worflow, err := spider.InitDefaultWorkflow(ctx)

	if err != nil {
		panic(err)
	}

	storage := worflow.Storage()

	uc := usecase.NewUsecase(storage)
	handler := apis.NewHandler(uc)

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return nil
	})

	// flows
	app.Get("/tenants/:tenant_id/flows", handler.ListFlows)
	app.Get("/tenants/:tenant_id/flows/:id", handler.GetFlow)
	app.Post("/tenants/:tenant_id/flows", handler.CreateFlow)
	app.Put("/tenants/:tenant_id/flows/:flow_id", handler.UpdateFlow)
	app.Delete("/tenants/:tenant_id/flows/:flow_id", handler.DeleteFlow)

	// actions
	app.Post("/tenants/:tenant_id/workflows/:workflow_id/actions/:key/disable", handler.DisableAction)
	app.Put("/tenants/:tenant_id/workflows/:workflow_id/actions/:key", handler.UpdateAction)

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
