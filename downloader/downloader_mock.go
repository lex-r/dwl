// Code generated by MockGen. DO NOT EDIT.
// Source: ./downloader/downloader.go

// Package downloader is a generated GoMock package.
package downloader

import (
	io "io"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHttpClient is a mock of HttpClient interface.
type MockHttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockHttpClientMockRecorder
}

// MockHttpClientMockRecorder is the mock recorder for MockHttpClient.
type MockHttpClientMockRecorder struct {
	mock *MockHttpClient
}

// NewMockHttpClient creates a new mock instance.
func NewMockHttpClient(ctrl *gomock.Controller) *MockHttpClient {
	mock := &MockHttpClient{ctrl: ctrl}
	mock.recorder = &MockHttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpClient) EXPECT() *MockHttpClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHttpClientMockRecorder) Do(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHttpClient)(nil).Do), req)
}

// MockSaver is a mock of Saver interface.
type MockSaver struct {
	ctrl     *gomock.Controller
	recorder *MockSaverMockRecorder
}

// MockSaverMockRecorder is the mock recorder for MockSaver.
type MockSaverMockRecorder struct {
	mock *MockSaver
}

// NewMockSaver creates a new mock instance.
func NewMockSaver(ctrl *gomock.Controller) *MockSaver {
	mock := &MockSaver{ctrl: ctrl}
	mock.recorder = &MockSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSaver) EXPECT() *MockSaverMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockSaver) Save(link string, r io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", link, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSaverMockRecorder) Save(link, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSaver)(nil).Save), link, r)
}