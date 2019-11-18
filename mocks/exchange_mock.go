// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pmdcosta/exchange/internal/handler (interfaces: Exchange)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	exchange "github.com/pmdcosta/exchange/pkg/exchange"
	currency "golang.org/x/text/currency"
	reflect "reflect"
)

// MockExchange is a mock of Exchange interface
type MockExchange struct {
	ctrl     *gomock.Controller
	recorder *MockExchangeMockRecorder
}

// MockExchangeMockRecorder is the mock recorder for MockExchange
type MockExchangeMockRecorder struct {
	mock *MockExchange
}

// NewMockExchange creates a new mock instance
func NewMockExchange(ctrl *gomock.Controller) *MockExchange {
	mock := &MockExchange{ctrl: ctrl}
	mock.recorder = &MockExchangeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExchange) EXPECT() *MockExchangeMockRecorder {
	return m.recorder
}

// GetRangeWithParams mocks base method
func (m *MockExchange) GetRangeWithParams(arg0 *exchange.Params, arg1, arg2 string) (map[string]map[currency.Unit]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRangeWithParams", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]map[currency.Unit]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRangeWithParams indicates an expected call of GetRangeWithParams
func (mr *MockExchangeMockRecorder) GetRangeWithParams(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRangeWithParams", reflect.TypeOf((*MockExchange)(nil).GetRangeWithParams), arg0, arg1, arg2)
}