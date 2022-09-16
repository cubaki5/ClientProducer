package consumer_adapter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/gommon/log"

	"clientProducer/internal/domain"
)

//go:generate  mockgen -source=adapter.go -destination=mocks/mock_consumer.go -package=mocks

type httpClient interface {
	PostReq(items []domain.Item) (*http.Response, error)
	GetReq() (*http.Response, error)
}

type ClientAdapter struct {
	client httpClient
}

func NewClientAdapter(httpC httpClient) *ClientAdapter {
	return &ClientAdapter{
		client: httpC,
	}
}

func (ca *ClientAdapter) GetBufferFreeSpace() (bufferFreeSpace int, err error) {
	resp, err := ca.client.GetReq()
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, errors.New("server is full")
	}

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bufferFreeSpace, err = strconv.Atoi(string(by))
	if err != nil {
		return 0, err
	}

	return bufferFreeSpace, nil
}

func (ca *ClientAdapter) PostBatch(batch []domain.Item) error {
	resp, err := ca.client.PostReq(batch)
	if err != nil {
		return err
	}

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("%s batch len %d", string(by), len(batch)))

	return nil
}
