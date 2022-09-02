package consumer

import (
	"bytes"
	"clientProducer/models"
	"clientProducer/module"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Consumer struct {
	ch               chan models.Batch
	consumerCapacity int
}

func NewConsumer() *Consumer {
	return &Consumer{
		ch:               make(chan models.Batch, 10000000),
		consumerCapacity: module.ConsumerBuffer,
	}
}

func (c *Consumer) Add(ctx context.Context, batch models.Batch) error {
	select {
	case c.ch <- batch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Consumer) Run() {
	var dividedBatch []models.Batch
	for batch := range c.ch {
		dividedBatch = c.batchDivider(batch)
		for _, batchSlice := range dividedBatch {
			by, err := postBatch(batchSlice)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(by))
			time.Sleep(module.ItemServeTime * 2)
		}
	}
}

func (c *Consumer) batchDivider(batch models.Batch) (dividedBatch []models.Batch) {
	for leftBound := 0; leftBound < len(batch); leftBound = leftBound + c.consumerCapacity {
		if leftBound+c.consumerCapacity > len(batch) {
			dividedBatch = append(dividedBatch, batch[leftBound:len(batch)])
		} else {
			dividedBatch = append(dividedBatch, batch[leftBound:len(batch)])
		}
	}
	return dividedBatch
}

func postBatch(batch models.Batch) (by []byte, err error) {
	resp, err := marshalingAndPost(batch)
	if err != nil {
		return nil, err
	}

	by, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return by, nil
}

func marshalingAndPost(batch models.Batch) (resp *http.Response, err error) {
	b, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	resp, err = http.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
