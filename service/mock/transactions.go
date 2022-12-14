// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Drozd0f/csv-app/service (interfaces: IFileHeader,IFile)

// Package mock_service is a generated GoMock package.
package mock_service

import (
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIFileHeader is a mock of IFileHeader interface.
type MockIFileHeader struct {
	ctrl     *gomock.Controller
	recorder *MockIFileHeaderMockRecorder
}

// MockIFileHeaderMockRecorder is the mock recorder for MockIFileHeader.
type MockIFileHeaderMockRecorder struct {
	mock *MockIFileHeader
}

// NewMockIFileHeader creates a new mock instance.
func NewMockIFileHeader(ctrl *gomock.Controller) *MockIFileHeader {
	mock := &MockIFileHeader{ctrl: ctrl}
	mock.recorder = &MockIFileHeaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFileHeader) EXPECT() *MockIFileHeaderMockRecorder {
	return m.recorder
}

// Open mocks base method.
func (m *MockIFileHeader) Open() (multipart.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(multipart.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockIFileHeaderMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockIFileHeader)(nil).Open))
}

// MockIFile is a mock of IFile interface.
type MockIFile struct {
	ctrl     *gomock.Controller
	recorder *MockIFileMockRecorder
}

// MockIFileMockRecorder is the mock recorder for MockIFile.
type MockIFileMockRecorder struct {
	mock *MockIFile
}

// NewMockIFile creates a new mock instance.
func NewMockIFile(ctrl *gomock.Controller) *MockIFile {
	mock := &MockIFile{ctrl: ctrl}
	mock.recorder = &MockIFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFile) EXPECT() *MockIFileMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockIFile) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockIFileMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIFile)(nil).Close))
}

// Read mocks base method.
func (m *MockIFile) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIFileMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIFile)(nil).Read), arg0)
}

// ReadAt mocks base method.
func (m *MockIFile) ReadAt(arg0 []byte, arg1 int64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAt", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAt indicates an expected call of ReadAt.
func (mr *MockIFileMockRecorder) ReadAt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAt", reflect.TypeOf((*MockIFile)(nil).ReadAt), arg0, arg1)
}

// Seek mocks base method.
func (m *MockIFile) Seek(arg0 int64, arg1 int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *MockIFileMockRecorder) Seek(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockIFile)(nil).Seek), arg0, arg1)
}
