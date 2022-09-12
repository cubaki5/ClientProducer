package consumer_adapter

import (
	"bytes"
	"clientProducer/internal/adapter/consumer_adapter/mocks"
	"clientProducer/internal/domain"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func initMock(t *testing.T) *mocks.MockhttpClient {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockhttpClient(ctrl)
	return mockCl
}

func newHttpResp(sc int, b string) *http.Response {
	return &http.Response{
		StatusCode: sc,
		Body:       ioutil.NopCloser(bytes.NewBufferString(b)),
	}
}

func TestClientAdapter(t *testing.T) {
	t.Run("Test Get request", func(t *testing.T) {
		getTests := []struct {
			testName     string
			statusCode   int
			respError    error
			expError     error
			freeSpace    string
			expFreeSpace int
		}{
			{
				testName:     "When response has status kod 400 and no error",
				statusCode:   http.StatusBadRequest,
				respError:    nil,
				expError:     errors.New("server is full"),
				freeSpace:    "0",
				expFreeSpace: 0,
			},
			{
				testName:     "When response has status kod 502 and  error",
				statusCode:   http.StatusBadGateway,
				respError:    errors.New("test error"),
				expError:     errors.New("test error"),
				freeSpace:    "0",
				expFreeSpace: 0,
			},
			{
				testName:     "When response has incorrect number of buffer free space",
				statusCode:   http.StatusOK,
				respError:    nil,
				expError:     errors.New("strconv.Atoi: parsing \"O\": invalid syntax"),
				freeSpace:    "O",
				expFreeSpace: 0,
			},
		}
		for _, getTest := range getTests {
			t.Run(getTest.testName, func(t *testing.T) {
				mockCl := initMock(t)
				resp := newHttpResp(getTest.statusCode, getTest.freeSpace)
				mockCl.EXPECT().GetReq().Return(resp, getTest.respError)
				ac := NewClientAdapter(mockCl)
				actVal, err := ac.GetBufferFreeSpace()

				require.EqualError(t, err, getTest.expError.Error())
				require.Equal(t, getTest.expFreeSpace, actVal)
			})
		}
		t.Run("When response has status kod 200 and no error", func(t *testing.T) {
			mockCl := initMock(t)
			resp := newHttpResp(http.StatusOK, "5")
			mockCl.EXPECT().GetReq().Return(resp, nil)
			ac := NewClientAdapter(mockCl)
			actVal, err := ac.GetBufferFreeSpace()

			require.NoError(t, err)
			require.Equal(t, 5, actVal)
		})

	})
	t.Run("Test Post Request", func(t *testing.T) {
		t.Run("When response has no error", func(t *testing.T) {
			mockCl := initMock(t)
			resp := newHttpResp(http.StatusOK, "test")
			mockCl.EXPECT().PostReq([]domain.Item{{}}).Return(resp, nil)
			ac := NewClientAdapter(mockCl)

			err := ac.PostBatch([]domain.Item{{}})

			require.NoError(t, err)
		})
		t.Run("When response has error", func(t *testing.T) {
			mockCl := initMock(t)
			resp := newHttpResp(http.StatusOK, "test")
			mockCl.EXPECT().PostReq(gomock.Any()).Return(resp, errors.New("test error"))
			ac := NewClientAdapter(mockCl)

			err := ac.PostBatch([]domain.Item{{}})

			require.EqualError(t, err, "test error")
		})

	})
}
