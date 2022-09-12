package consumer_queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"clientProducer/internal/domain"

	"github.com/labstack/gommon/log"
)

var (
	cons *consumerQueue
	once sync.Once
)

//go:generate mockgen --destination=mocks/mock_client_adapter.go --package=mocks

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
	dautCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	select {
	case c.ch <- batch:
		return nil
	case <-dautCtx.Done():
		err := errors.New("handlers queue full, check back later")
		log.Error(err)
		return err
	}
}

func (c *consumerQueue) Run() {
	for batch := range c.ch {

		var err error

		c.bufferFreeSpace, err = c.getFreeSpace()
		if err != nil {
			log.Error(err)
			break
		}

		batch = c.batchDivider(batch)

		err = c.cliCon.PostBatch(batch)
		if err != nil {
			log.Error(err)
		}
	}
}

func (c *consumerQueue) getFreeSpace() (bufFreeSpace int, err error) {
	ticker := time.NewTicker(itemServeTime)
	for range ticker.C {
		bufFreeSpace, err = c.cliCon.GetBufferFreeSpace()
		if err != nil {
			log.Error(err)
			return 0, err
		}
		if bufFreeSpace > 0 {
			break
		}
	}
	return bufFreeSpace, nil
}

func (c *consumerQueue) batchDivider(batch []domain.Item) (dividedBatch []domain.Item) {

	dividedBatch = batch
	if len(batch) > c.bufferFreeSpace {
		con, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := c.Add(con, batch[c.bufferFreeSpace:])
		if err != nil {
			log.Error(err)
		}
		dividedBatch = batch[0:c.bufferFreeSpace]
	}

	return dividedBatch
}
