package append_items

import (
	"clientProducer/internal/domain"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	Handle(items []domain.Item) error
}

type handlerEcho struct {
	h Handler
}

func NewAppendItems(h Handler) *handlerEcho {
	return &handlerEcho{h: h}
}

func (he *handlerEcho) Handle(echoCtx echo.Context) error {
	items, err := parseRequestBody(echoCtx)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}

	if err = he.h.Handle(items); err != nil {
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}

	return echoCtx.NoContent(http.StatusOK)
}

func parseRequestBody(echoCtx echo.Context) (items []domain.Item, err error) {
	by, err := ioutil.ReadAll(echoCtx.Request().Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(by, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
