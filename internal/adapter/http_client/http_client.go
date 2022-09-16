package http_client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"clientProducer/internal/domain"
)

const (
	serverURL = "http://localhost:1323"
)

type Client struct {
	client http.Client
}

func NewClientAdapter() *Client {
	return &Client{
		client: http.Client{},
	}
}

func (c *Client) PostReq(items []domain.Item) (*http.Response, error) {
	b, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(serverURL, "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetReq() (*http.Response, error) {
	resp, err := c.client.Get(serverURL + "/buffer")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
