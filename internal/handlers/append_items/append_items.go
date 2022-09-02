package append_items

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"clientProducer/internal/domain"

	"github.com/labstack/echo/v4"
)

type UseCase interface {
	Run(ctx context.Context, items []domain.Item) error
}

type appendItemsHandler struct {
	uc UseCase
}

func NewAppendItemHandler(uc UseCase) *appendItemsHandler {
	return &appendItemsHandler{uc: uc}
}

func (ai *appendItemsHandler) Handle(echoCtx echo.Context) error {
	items, err := ai.parseRequestBody(echoCtx)
	if err != nil {
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}
	con := context.Background()
	if err = ai.uc.Run(con, items); err != nil {
		return echoCtx.String(http.StatusBadRequest, err.Error())
	}
	return echoCtx.NoContent(http.StatusOK)
}

func (ai *appendItemsHandler) parseRequestBody(echoCtx echo.Context) (items []domain.Item, err error) {
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
