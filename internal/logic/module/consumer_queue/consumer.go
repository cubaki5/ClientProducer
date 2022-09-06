package consumer_queue

import (
	"clientProducer/internal/domain"
	"context"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

var (
	cons *consumerQueue
	once sync.Once
)

type ClientConsumer interface {
	GetBufferFreeSpace() (bufferFreeSpace int, err error)
	PostBatch(batch []domain.Item) error
}

type consumerQueue struct {
	ch              chan []domain.Item
	bufferFreeSpace int
	cliCon          ClientConsumer
}

func NewConsumerQueue(cliCon ClientConsumer) *consumerQueue {
	once.Do(func() {
		if cons == nil {
			cons = &consumerQueue{
				ch:              make(chan []domain.Item, maxBatch),
				bufferFreeSpace: consumerBuffer,
				cliCon:          cliCon,
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
	var err error
	for batch := range c.ch {
		for {
			c.bufferFreeSpace, err = c.cliCon.GetBufferFreeSpace()
			if err != nil {
				log.Error(err)
				time.Sleep(PanicDuration)
			} else {
				if c.bufferFreeSpace >= len(batch) {
					break
				} else {
					time.Sleep(itemServeTime)
				}
			}
		}
		err = c.cliCon.PostBatch(batch)
		if err != nil {
			log.Error(err)
		}
	}

}

func (c *consumerQueue) waiterForFreeSpace() {

}

func (c *consumerQueue) batchDivider(batch []domain.Item) (dividedBatch [][]domain.Item) {
	dividedBatch = append(dividedBatch, batch[0:c.bufferFreeSpace])
	dividedBatch = append(dividedBatch, batch[c.bufferFreeSpace:])
	return dividedBatch
}
