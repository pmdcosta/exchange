package exchange_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pmdcosta/exchange/mocks"
	"github.com/pmdcosta/exchange/pkg/exchange"
	"github.com/rs/zerolog"
)

// Client is a test wrapper
type Client struct {
	exchange.Client
	mockCtrl  *gomock.Controller
	mockCache *mocks.MockCache
}

func newClient(t *testing.T, address string) *Client {
	mockCtrl := gomock.NewController(t)
	mockCache := mocks.NewMockCache(mockCtrl)
	logger := zerolog.Nop()

	return &Client{
		Client:    *exchange.NewClient(&logger, exchange.SetAddress(address)),
		mockCtrl:  mockCtrl,
		mockCache: mockCache,
	}
}

func (s *Client) Finish() {
	s.mockCtrl.Finish()
}
