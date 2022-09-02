package adapter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"clientProducer/internal/domain"
)

func PostBatch(batch []domain.Item) (by []byte, err error) {
	resp, err := marshallingAndPost(batch)
	if err != nil {
		return nil, err
	}

	by, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return by, nil
}

func marshallingAndPost(batch []domain.Item) (resp *http.Response, err error) {
	b, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	resp, err = http.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
