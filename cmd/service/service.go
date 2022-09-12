package main

import (
	"clientProducer/internal/adapter/consumer_adapter"
	append_items_echo_handler "clientProducer/internal/adapter/handlers/append_items"
	httpclient "clientProducer/internal/adapter/httpclient"
	append_items_handler "clientProducer/internal/handlers/append_items"
	"clientProducer/internal/logic/module/consumer_queue"
	"clientProducer/internal/logic/usecase/append_items"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	//debug := flag.Bool("debug", true, "sets log level to debug")
	//
	//flag.Parse()
	//log.SetLevel(log.ERROR)
	//if *debug {
	//	log.SetLevel(log.DEBUG)
	//}
	log.SetLevel(log.DEBUG)
	consCli := consumer_adapter.NewClientAdapter(httpclient.NewClientAdapter())
	cons := consumer_queue.NewConsumerQueue(consCli)
	useCase := append_items.NewAppendItemsUseCase(cons)

	h := append_items_echo_handler.NewAppendItems(append_items_handler.NewAppendItemHandler(useCase))

	e := echo.New()
	e.POST("/", h.Handle)

	e.Logger.Fatal(e.Start(":1324"))
}
