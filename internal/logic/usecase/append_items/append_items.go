package append_items

import (
	"context"
	"github.com/labstack/gommon/log"

	"clientProducer/internal/domain"
)

//go:generate mockgen --destination=mocks/mock_consumer.go --package=mocks
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
	err := a.consumer.Add(ctx, items)
	if err != nil {
		log.Error(err)
	}
	return err
}
