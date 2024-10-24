// Code generated by MockGen. DO NOT EDIT.
// Source: parser_repository.go
//
// Generated by this command:
//
//	mockgen -source=parser_repository.go -destination=./mock/mock_parser_repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gofeed "github.com/mmcdole/gofeed"
	gomock "go.uber.org/mock/gomock"
)

// MockParserRepository is a mock of ParserRepository interface.
type MockParserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockParserRepositoryMockRecorder
}

// MockParserRepositoryMockRecorder is the mock recorder for MockParserRepository.
type MockParserRepositoryMockRecorder struct {
	mock *MockParserRepository
}

// NewMockParserRepository creates a new mock instance.
func NewMockParserRepository(ctrl *gomock.Controller) *MockParserRepository {
	mock := &MockParserRepository{ctrl: ctrl}
	mock.recorder = &MockParserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParserRepository) EXPECT() *MockParserRepositoryMockRecorder {
	return m.recorder
}

// ParseURLWithContext mocks base method.
func (m *MockParserRepository) ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseURLWithContext", url, ctx)
	ret0, _ := ret[0].(*gofeed.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseURLWithContext indicates an expected call of ParseURLWithContext.
func (mr *MockParserRepositoryMockRecorder) ParseURLWithContext(url, ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseURLWithContext", reflect.TypeOf((*MockParserRepository)(nil).ParseURLWithContext), url, ctx)
}
