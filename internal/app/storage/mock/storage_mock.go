// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/northmule/shorturl/internal/app/storage (interfaces: DBQuery,DBRow)
//
// Generated by this command:
//
//	mockgen -destination=internal/app/storage/mock/storage_mock.go -package=mocks github.com/northmule/shorturl/internal/app/storage DBQuery,DBRow
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDBQuery is a mock of DBQuery interface.
type MockDBQuery struct {
	ctrl     *gomock.Controller
	recorder *MockDBQueryMockRecorder
}

// MockDBQueryMockRecorder is the mock recorder for MockDBQuery.
type MockDBQueryMockRecorder struct {
	mock *MockDBQuery
}

// NewMockDBQuery creates a new mock instance.
func NewMockDBQuery(ctrl *gomock.Controller) *MockDBQuery {
	mock := &MockDBQuery{ctrl: ctrl}
	mock.recorder = &MockDBQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBQuery) EXPECT() *MockDBQueryMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockDBQuery) Begin() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockDBQueryMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockDBQuery)(nil).Begin))
}

// ExecContext mocks base method.
func (m *MockDBQuery) ExecContext(arg0 context.Context, arg1 string, arg2 ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockDBQueryMockRecorder) ExecContext(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockDBQuery)(nil).ExecContext), varargs...)
}

// PingContext mocks base method.
func (m *MockDBQuery) PingContext(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockDBQueryMockRecorder) PingContext(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockDBQuery)(nil).PingContext), arg0)
}

// QueryContext mocks base method.
func (m *MockDBQuery) QueryContext(arg0 context.Context, arg1 string, arg2 ...any) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockDBQueryMockRecorder) QueryContext(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockDBQuery)(nil).QueryContext), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockDBQuery) QueryRowContext(arg0 context.Context, arg1 string, arg2 ...any) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockDBQueryMockRecorder) QueryRowContext(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockDBQuery)(nil).QueryRowContext), varargs...)
}

// MockDBRow is a mock of DBRow interface.
type MockDBRow struct {
	ctrl     *gomock.Controller
	recorder *MockDBRowMockRecorder
}

// MockDBRowMockRecorder is the mock recorder for MockDBRow.
type MockDBRowMockRecorder struct {
	mock *MockDBRow
}

// NewMockDBRow creates a new mock instance.
func NewMockDBRow(ctrl *gomock.Controller) *MockDBRow {
	mock := &MockDBRow{ctrl: ctrl}
	mock.recorder = &MockDBRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBRow) EXPECT() *MockDBRowMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockDBRow) Scan(arg0 ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockDBRowMockRecorder) Scan(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockDBRow)(nil).Scan), arg0...)
}
