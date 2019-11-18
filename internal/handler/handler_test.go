package handler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pmdcosta/exchange/internal/handler"
	"github.com/pmdcosta/exchange/mocks"
	"github.com/rs/zerolog"
)

// Handler is a test handler wrapper
type Handler struct {
	*handler.Handler
	mockCtrl *gomock.Controller
	exchange *mocks.MockExchange
}

// New returns a new test http handler
func New(t *testing.T, recommendation handler.Recommendation) *Handler {
	mockCtrl := gomock.NewController(t)
	mockExchange := mocks.NewMockExchange(mockCtrl)
	logger := zerolog.Nop()

	return &Handler{
		Handler:  handler.New(&logger, mockExchange, recommendation),
		mockCtrl: mockCtrl,
		exchange: mockExchange,
	}
}

// Finish finishes the test
func (s *Handler) Finish() {
	s.mockCtrl.Finish()
}
