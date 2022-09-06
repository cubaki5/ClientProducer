package append_items

import (
	"clientProducer/internal/domain"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockConsumer struct{}

func (mc MockConsumer) Add(ctx context.Context, items []domain.Item) error {
	if items == nil {
		return fmt.Errorf("h")
	}
	return nil
}

func TestAppendItems_Run(t *testing.T) {
	usecase := NewAppendItemsUseCase(MockConsumer{})
	err := usecase.Run(context.Background(), []domain.Item{{}, {}})
	assert.NoError(t, err)
}
