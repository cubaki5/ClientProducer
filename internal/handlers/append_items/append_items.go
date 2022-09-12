package append_items

import (
	"clientProducer/internal/domain"
	"context"
)

type UseCase interface {
	Run(ctx context.Context, items []domain.Item) error
}

type appendItemsHandler struct {
	uc UseCase
}

func NewAppendItemHandler(uc UseCase) *appendItemsHandler {
	return &appendItemsHandler{uc: uc}
}

func (ai *appendItemsHandler) Handle(items []domain.Item) error {
	return ai.uc.Run(context.Background(), items)
}
