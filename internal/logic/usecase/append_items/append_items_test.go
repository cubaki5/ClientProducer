package append_items

import (
	"context"
	"errors"
	"testing"

	"clientProducer/internal/domain"
	mock_append_items "clientProducer/internal/logic/usecase/append_items/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initMock(t *testing.T) *mock_append_items.MockConsumer {
	ctrl := gomock.NewController(t)
	mockedConsumer := mock_append_items.NewMockConsumer(ctrl)
	return mockedConsumer
}

func TestNewAppendItemsUseCase(t *testing.T) {
	t.Run("Consumer gets the correct slice of items", func(t *testing.T) {
		mockedConsumer := initMock(t)
		mockedConsumer.EXPECT().Add(gomock.Any(), []domain.Item{{}})
		useCase := NewAppendItemsUseCase(mockedConsumer)

		_ = useCase.Run(context.Background(), []domain.Item{{}})
	})
	t.Run("UseCase return correct err", func(t *testing.T) {
		mockedConsumer := initMock(t)
		mockedConsumer.EXPECT().Add(gomock.Any(), gomock.Any()).Return(errors.New("test error"))
		useCase := NewAppendItemsUseCase(mockedConsumer)

		err := useCase.Run(context.Background(), []domain.Item{})

		assert.EqualError(t, err, "test error")
	})
	t.Run("Consumer gets correct context", func(t *testing.T) {
		mockedConsumer := initMock(t)
		mockedConsumer.EXPECT().Add(context.Background(), gomock.Any())
		useCase := NewAppendItemsUseCase(mockedConsumer)

		_ = useCase.Run(context.Background(), []domain.Item{})
	})
}
