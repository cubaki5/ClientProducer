package adapter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"clientProducer/internal/domain"
)

type Client interface {
	GetBufferFreeSpace() (bufferFreeSpace int, err error)
}

type Adapter struct {
	client http.Client
}

//func NewClient() *client {
//	return &client{
//		Transport:     nil,
//		CheckRedirect: nil,
//		Jar:           nil,
//		Timeout:       0,
//	}
//}

func GetBufferFreeSpace() (bufferFreeSpace int, err error) {
	client.Get()
}

func PostBatch(batch []domain.Item) (by []byte, err error) {
	http.New
	resp, err := marshallingAndPostReq(batch)
	if err != nil {
		return nil, err
	}
	by, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return by, nil
}

func getReq() (resp *http.Response, err error) {

}

func (a *Adapter) marshallingAndPostReq(batch []domain.Item) (resp *http.Response, err error) {
	b, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	a.client.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	resp, err = http.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
