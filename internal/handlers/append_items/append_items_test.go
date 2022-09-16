package append_items

import (
	"errors"
	"testing"

	"clientProducer/internal/domain"
	"clientProducer/internal/handlers/append_items/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func initMock(t *testing.T) *mocks.MockUseCase {
	ctrl := gomock.NewController(t)
	mockUseCase := mocks.NewMockUseCase(ctrl)
	return mockUseCase
}

func TestAppendItemsHandler_Handle(t *testing.T) {
	t.Run("When send items without error", func(t *testing.T) {
		mockUseCase := initMock(t)
		handler := NewAppendItemHandler(mockUseCase)
		mockUseCase.EXPECT().Run(gomock.Any(), []domain.Item{{}, {}, {}}).Return(nil)

		err := handler.Handle([]domain.Item{{}, {}, {}})
		require.NoError(t, err)
	})
	t.Run("When send items with error", func(t *testing.T) {
		mockUseCase := initMock(t)
		handler := NewAppendItemHandler(mockUseCase)
		mockUseCase.EXPECT().Run(gomock.Any(), []domain.Item{{}, {}, {}}).Return(errors.New("test error"))

		err := handler.Handle([]domain.Item{{}, {}, {}})
		require.EqualError(t, err, "test error")
	})
}
