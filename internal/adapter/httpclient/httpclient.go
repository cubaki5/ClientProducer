package httpclient

import (
	"bytes"
	"clientProducer/internal/domain"
	"encoding/json"
	"net/http"
)

const (
	serverUrl = "http://localhost:1323"
)

type Client struct {
	client http.Client
}

func NewClientAdapter() *Client {
	return &Client{
		client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
	}
}

func (c *Client) PostReq(items []domain.Item) (*http.Response, error) {
	b, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(serverUrl, "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetReq() (*http.Response, error) {
	resp, err := c.client.Get(serverUrl + "/buffer")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
