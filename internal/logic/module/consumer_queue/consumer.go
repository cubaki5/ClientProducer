package consumer_queue

import (
	"context"
	"sync"
	"time"

	"clientProducer/internal/adapter"
	"clientProducer/internal/domain"

	"github.com/labstack/gommon/log"
)

var (
	cons *consumerQueue
	once sync.Once
)

type consumerQueue struct {
	ch               chan []domain.Item
	consumerCapacity int
}

func NewConsumerQueue() *consumerQueue {
	once.Do(func() {
		if cons == nil {
			cons = &consumerQueue{
				ch:               make(chan []domain.Item, maxBatch),
				consumerCapacity: consumerBuffer,
			}
			go cons.Run()
		}
	})
	return cons
}

func (c *consumerQueue) Add(ctx context.Context, batch []domain.Item) error {
	select {
	case c.ch <- batch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *consumerQueue) Run() {
	var dividedBatch [][]domain.Item
	for batch := range c.ch {
		dividedBatch = c.batchDivider(batch)
		for _, batchSlice := range dividedBatch {
			by, err := adapter.PostBatch(batchSlice)
			if err != nil {
				log.Error(err)
			}
			log.Debug(string(by))
			time.Sleep(itemServeTime * 2)
		}
	}
}

func (c *consumerQueue) batchDivider(batch []domain.Item) (dividedBatch [][]domain.Item) {
	for leftBound := 0; leftBound < len(batch); leftBound = leftBound + c.consumerCapacity {
		if leftBound+c.consumerCapacity > len(batch) {
			dividedBatch = append(dividedBatch, batch[leftBound:])
		} else {
			dividedBatch = append(dividedBatch, batch[leftBound:leftBound+c.consumerCapacity])
		}
	}
	return dividedBatch
}
