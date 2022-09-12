package consumer_queue

import (
	"clientProducer/internal/domain"
	"clientProducer/internal/logic/module/consumer_queue/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func initMock(t *testing.T) *mocks.MockClientConsumer {
	ctrl := gomock.NewController(t)
	mockCliCon := mocks.NewMockClientConsumer(ctrl)
	return mockCliCon
}

func TestConsumerQueue(t *testing.T) {

	t.Run("Test Add function", func(t *testing.T) {
		t.Run("Add function return err after time out", func(t *testing.T) {
			mockCliCon := initMock(t)
			conQueue := &consumerQueue{
				ch:     make(chan []domain.Item, maxBatch),
				cliCon: mockCliCon,
			}
			consChan := make(chan []domain.Item, maxBatch)
			for i := 0; i < maxBatch; i++ {
				consChan <- []domain.Item{{}}
			}
			conQueue.ch = consChan

			err := conQueue.Add(context.Background(), []domain.Item{{}})

			require.EqualError(t, err, "handlers queue full, check back later")
		})
		t.Run("Add function return nil after adding items", func(t *testing.T) {
			mockCliCon := initMock(t)
			conQueue := &consumerQueue{
				ch:     make(chan []domain.Item, maxBatch),
				cliCon: mockCliCon,
			}

			err := conQueue.Add(context.Background(), []domain.Item{{}})

			require.NoError(t, err)
		})
	})

	t.Run("Test getFreeSpace function", func(t *testing.T) {
		t.Run("When GetBufferFreeSpace returns err", func(t *testing.T) {
			mockCliCon := initMock(t)
			conQueue := &consumerQueue{
				ch:              make(chan []domain.Item, maxBatch),
				bufferFreeSpace: consumerBuffer,
				cliCon:          mockCliCon,
			}
			mockCliCon.EXPECT().GetBufferFreeSpace().Return(0, errors.New("test error"))

			freeSpace, err := conQueue.getFreeSpace()
			require.EqualError(t, err, "test error")
			require.Equal(t, 0, freeSpace)
		})
		t.Run("When GetBufferFreeSpace returns freeSpace", func(t *testing.T) {
			mockCliCon := initMock(t)
			conQueue := &consumerQueue{
				ch:              make(chan []domain.Item, maxBatch),
				bufferFreeSpace: consumerBuffer,
				cliCon:          mockCliCon,
			}
			mockCliCon.EXPECT().GetBufferFreeSpace().Return(5, nil)

			freeSpace, _ := conQueue.getFreeSpace()
			require.Equal(t, 5, freeSpace)
		})
	})

	t.Run("Test BatchDivider function", func(t *testing.T) {
		mockCliCon := initMock(t)
		conQueue := &consumerQueue{
			ch:              make(chan []domain.Item, maxBatch),
			bufferFreeSpace: consumerBuffer,
			cliCon:          mockCliCon,
		}

		dividerTests := []struct {
			testName string
			batch    []domain.Item
			expBatch []domain.Item
		}{
			{
				testName: "When len of batch is equal to bufferFreeSpace",
				batch:    make([]domain.Item, conQueue.bufferFreeSpace),
				expBatch: make([]domain.Item, conQueue.bufferFreeSpace),
			},
			{
				testName: "When len of batch is twice more than bufferFreeSpace",
				batch:    make([]domain.Item, 2*conQueue.bufferFreeSpace),
				expBatch: make([]domain.Item, conQueue.bufferFreeSpace),
			},
			{
				testName: "When len of batch is less than bufferFreeSpace",
				batch:    make([]domain.Item, conQueue.bufferFreeSpace-1),
				expBatch: make([]domain.Item, conQueue.bufferFreeSpace-1),
			},
			{
				testName: "When bufferFreeSpace is null",
				batch:    make([]domain.Item, 0),
				expBatch: []domain.Item{},
			},
		}

		for _, dividerTest := range dividerTests {
			t.Run(dividerTest.testName, func(t *testing.T) {
				actBatch := conQueue.batchDivider(dividerTest.batch)

				require.Equal(t, dividerTest.expBatch, actBatch)
			})
		}

	})

	t.Run("Test Run function", func(t *testing.T) {
		t.Run("Test channel release", func(t *testing.T) {
			mockCliCon := initMock(t)
			mockCliCon.EXPECT().GetBufferFreeSpace().Return(5, nil).AnyTimes()
			mockCliCon.EXPECT().PostBatch(gomock.Any()).Return(nil).AnyTimes()

			conQueue := &consumerQueue{
				ch:              make(chan []domain.Item, maxBatch),
				bufferFreeSpace: consumerBuffer,
				cliCon:          mockCliCon,
			}
			conQueue.ch <- []domain.Item{{}, {}, {}, {}, {}}

			go conQueue.Run()
			time.Sleep(time.Second)

			require.Equal(t, 0, len(conQueue.ch))
		})
		t.Run("PostBatch get correct batch", func(t *testing.T) {
			mockCliCon := initMock(t)
			mockCliCon.EXPECT().GetBufferFreeSpace().Return(5, nil).AnyTimes()
			mockCliCon.EXPECT().PostBatch([]domain.Item{{}, {}, {}, {}, {}}).Return(nil).AnyTimes()
			conQueue := &consumerQueue{
				ch:              make(chan []domain.Item, maxBatch),
				bufferFreeSpace: consumerBuffer,
				cliCon:          mockCliCon,
			}

			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			err := conQueue.Add(ctx, []domain.Item{{}, {}, {}, {}, {}})
			require.NoError(t, err)

			go conQueue.Run()
			time.Sleep(time.Second * 2)
		})
	})
}
