package append_items

import (
	"context"

	"clientProducer/internal/domain"
)

type Consumer interface {
	Add(ctx context.Context, items []domain.Item) error
}

type appendItems struct {
	consumer Consumer
}

func NewAppendItemsUseCase(cons Consumer) *appendItems {
	return &appendItems{consumer: cons}
}

func (a appendItems) Run(ctx context.Context, items []domain.Item) error {
	if err := a.consumer.Add(ctx, items); err != nil {
		return err
	}
	return nil
}
