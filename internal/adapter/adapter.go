package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"clientProducer/internal/domain"

	"github.com/labstack/gommon/log"
)

type ClientAdapter struct {
	client http.Client
}

func NewClientAdapter() *ClientAdapter {
	return &ClientAdapter{
		client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
	}
}

func (ca *ClientAdapter) GetBufferFreeSpace() (bufferFreeSpace int, err error) {
	resp, err := ca.client.Get("http://localhost:1323/buffer")
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
	resp, err := ca.marshallingAndPostReq(batch)
	if err != nil {
		return err
	}

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debug(string(by), " batch len ", len(batch))

	return nil
}

func (ca *ClientAdapter) marshallingAndPostReq(batch []domain.Item) (resp *http.Response, err error) {
	b, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	resp, err = ca.client.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
