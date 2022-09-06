package consumer_queue

import "time"

const (
	consumerBuffer int           = 5
	itemServeTime                = 2 * time.Second
	PanicDuration  time.Duration = 30 * time.Second
	maxBatch                     = 1000000
)
