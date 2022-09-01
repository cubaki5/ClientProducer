package main

import (
	"bytes"
	"clientProducer/models"
	"clientProducer/module"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Producer struct {
	batch models.Batch
}

func newProducer(numOfItem int) (producer Producer) {
	return Producer{
		batch: make([]models.Item, numOfItem),
	}
}

func main() {
	var numsOfItem []int = []int{20, 16, 4, 4, 4, 4}

	buf := NewBuffer()
	for _, numOfItem := range numsOfItem {
		buf.AddBuffer(newProducer(numOfItem).batch)
		//time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
	}

	buf.ReleaseBuffer()
}

type Buffer struct {
	dB             map[int]models.Batch
	x              int
	numOfFreeSpace int
	dBLock         sync.Mutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		dB:             make(map[int]models.Batch),
		x:              1,
		numOfFreeSpace: module.ConsumerBuffer,
	}
}

func (b *Buffer) AddBuffer(batch models.Batch) {
	b.dBLock.Lock()
	defer b.dBLock.Unlock()
	b.dB[b.x] = batch
	b.x++
}

func (b *Buffer) CountPostBatch(batch models.Batch, initInd, termInd int) (by []byte) {
	b.numOfFreeSpace -= len(batch[initInd:termInd])

	resp, err := PostBatch(batch[initInd:termInd])
	if err != nil {
		log.Println(err)
	}

	by, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(module.ItemServeTime * 2)
	b.numOfFreeSpace += len(batch[initInd:termInd])
	return by
}

func (b *Buffer) ReleaseBuffer() {
	for _, batch := range b.dB {
		if len(batch) <= b.numOfFreeSpace {
			log.Println("начал отправку1")
			by := b.CountPostBatch(batch, 0, len(batch))
			log.Println(string(by))
		} else {
			for i := 0; i < len(batch); i = i + b.numOfFreeSpace {
				if i+b.numOfFreeSpace > len(batch) {
					log.Println("начал отправку2")
					by := b.CountPostBatch(batch, i, len(batch))
					log.Println(string(by))
				} else {
					log.Println("начал отправку3")
					by := b.CountPostBatch(batch, i, i+b.numOfFreeSpace)
					log.Println(string(by))
				}
			}
		}
	}
}

func PostBatch(batch models.Batch) (resp *http.Response, err error) {
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
