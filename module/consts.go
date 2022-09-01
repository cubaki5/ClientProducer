package module

import "time"

const (
	ConsumerBuffer int           = 10
	ItemServeTime  time.Duration = 2 * time.Second
	PanicDuration  time.Duration = 10 * time.Second
)
