package usecase

import (
	"context"

	"github.com/targc/spider-go/pkg/spider"
)

type Usecase struct {
	storage spider.WorkflowStorageAdapter
	ctx     context.Context
}

func NewUsecase(storage spider.WorkflowStorageAdapter, ctx context.Context) *Usecase {
	return &Usecase{
		storage: storage,
		ctx:     ctx,
	}
}
