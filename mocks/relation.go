// Code generated by MockGen. DO NOT EDIT.
// Source: models/relation.go

// Package mock_models is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRelationDatabase is a mock of RelationDatabase interface
type MockRelationDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockRelationDatabaseMockRecorder
}

// MockRelationDatabaseMockRecorder is the mock recorder for MockRelationDatabase
type MockRelationDatabaseMockRecorder struct {
	mock *MockRelationDatabase
}

// NewMockRelationDatabase creates a new mock instance
func NewMockRelationDatabase(ctrl *gomock.Controller) *MockRelationDatabase {
	mock := &MockRelationDatabase{ctrl: ctrl}
	mock.recorder = &MockRelationDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRelationDatabase) EXPECT() *MockRelationDatabaseMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockRelationDatabase) Find(arg0, arg1, arg2 string) ([][]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1, arg2)
	ret0, _ := ret[0].([][]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockRelationDatabaseMockRecorder) Find(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockRelationDatabase)(nil).Find), arg0, arg1, arg2)
}

// Connect mocks base method
func (m *MockRelationDatabase) Connect(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockRelationDatabaseMockRecorder) Connect(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockRelationDatabase)(nil).Connect), arg0, arg1, arg2)
}

// Disconnect mocks base method
func (m *MockRelationDatabase) Disconnect(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect
func (mr *MockRelationDatabaseMockRecorder) Disconnect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockRelationDatabase)(nil).Disconnect), arg0, arg1)
}

// Following mocks base method
func (m *MockRelationDatabase) Following(arg0 string) ([][]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Following", arg0)
	ret0, _ := ret[0].([][]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Following indicates an expected call of Following
func (mr *MockRelationDatabaseMockRecorder) Following(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Following", reflect.TypeOf((*MockRelationDatabase)(nil).Following), arg0)
}

// Followers mocks base method
func (m *MockRelationDatabase) Followers(arg0 string) ([][]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Followers", arg0)
	ret0, _ := ret[0].([][]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Followers indicates an expected call of Followers
func (mr *MockRelationDatabaseMockRecorder) Followers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Followers", reflect.TypeOf((*MockRelationDatabase)(nil).Followers), arg0)
}
