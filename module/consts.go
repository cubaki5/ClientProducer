package module

import "time"

const (
	ConsumerBuffer int           = 10
	ItemServeTime  time.Duration = 2 * time.Millisecond
	PanicDuration  time.Duration = 10 * time.Second
)
