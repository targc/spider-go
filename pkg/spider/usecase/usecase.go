package usecase

import (
	"github.com/targc/spider-go/pkg/spider"
)

type Usecase struct {
	storage spider.WorkflowStorageAdapter
}

func NewUsecase(storage spider.WorkflowStorageAdapter) *Usecase {
	return &Usecase{
		storage: storage,
	}
}
