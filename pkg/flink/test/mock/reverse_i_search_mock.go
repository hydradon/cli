// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/confluentinc/cli/v3/pkg/flink/internal/reverseisearch (interfaces: ReverseISearch)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockReverseISearch is a mock of ReverseISearch interface.
type MockReverseISearch struct {
	ctrl     *gomock.Controller
	recorder *MockReverseISearchMockRecorder
}

// MockReverseISearchMockRecorder is the mock recorder for MockReverseISearch.
type MockReverseISearchMockRecorder struct {
	mock *MockReverseISearch
}

// NewMockReverseISearch creates a new mock instance.
func NewMockReverseISearch(ctrl *gomock.Controller) *MockReverseISearch {
	mock := &MockReverseISearch{ctrl: ctrl}
	mock.recorder = &MockReverseISearchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReverseISearch) EXPECT() *MockReverseISearchMockRecorder {
	return m.recorder
}

// ReverseISearch mocks base method.
func (m *MockReverseISearch) ReverseISearch(arg0 []string, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReverseISearch", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// ReverseISearch indicates an expected call of ReverseISearch.
func (mr *MockReverseISearchMockRecorder) ReverseISearch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReverseISearch", reflect.TypeOf((*MockReverseISearch)(nil).ReverseISearch), arg0, arg1)
}
