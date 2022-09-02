package main

import (
	"clientProducer/models"
	"clientProducer/module/consumer"
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	await := make(chan struct{})

	cons := consumer.NewConsumer()

	go cons.Run()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	for i := 0; i < 10; i++ {
		batch := make([]models.Item, rand.Int31n(100))
		err := cons.Add(ctx, batch)
		if err != nil {
			log.Println(err)
		}
	}

	await <- struct{}{}
}
