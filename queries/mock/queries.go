// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kanatsanan6/hrm/queries (interfaces: Queries)

// Package mock_queries is a generated GoMock package.
package mock_queries

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/kanatsanan6/hrm/model"
	queries "github.com/kanatsanan6/hrm/queries"
)

// MockQueries is a mock of Queries interface.
type MockQueries struct {
	ctrl     *gomock.Controller
	recorder *MockQueriesMockRecorder
}

// MockQueriesMockRecorder is the mock recorder for MockQueries.
type MockQueriesMockRecorder struct {
	mock *MockQueries
}

// NewMockQueries creates a new mock instance.
func NewMockQueries(ctrl *gomock.Controller) *MockQueries {
	mock := &MockQueries{ctrl: ctrl}
	mock.recorder = &MockQueriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueries) EXPECT() *MockQueriesMockRecorder {
	return m.recorder
}

// CreateCompany mocks base method.
func (m *MockQueries) CreateCompany(arg0 queries.CreateCompanyArgs) (model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", arg0)
	ret0, _ := ret[0].(model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockQueriesMockRecorder) CreateCompany(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockQueries)(nil).CreateCompany), arg0)
}

// CreateLeave mocks base method.
func (m *MockQueries) CreateLeave(arg0 queries.CreateLeaveArgs) (model.Leave, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLeave", arg0)
	ret0, _ := ret[0].(model.Leave)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLeave indicates an expected call of CreateLeave.
func (mr *MockQueriesMockRecorder) CreateLeave(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLeave", reflect.TypeOf((*MockQueries)(nil).CreateLeave), arg0)
}

// CreateUser mocks base method.
func (m *MockQueries) CreateUser(arg0 queries.CreateUserArgs) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockQueriesMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockQueries)(nil).CreateUser), arg0)
}

// DeleteUser mocks base method.
func (m *MockQueries) DeleteUser(arg0 model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockQueriesMockRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockQueries)(nil).DeleteUser), arg0)
}

// FindCompanyByID mocks base method.
func (m *MockQueries) FindCompanyByID(arg0 uint) (model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCompanyByID", arg0)
	ret0, _ := ret[0].(model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCompanyByID indicates an expected call of FindCompanyByID.
func (mr *MockQueriesMockRecorder) FindCompanyByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCompanyByID", reflect.TypeOf((*MockQueries)(nil).FindCompanyByID), arg0)
}

// FindUserByEmail mocks base method.
func (m *MockQueries) FindUserByEmail(arg0 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockQueriesMockRecorder) FindUserByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockQueries)(nil).FindUserByEmail), arg0)
}

// FindUserByForgetPasswordToken mocks base method.
func (m *MockQueries) FindUserByForgetPasswordToken(arg0 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByForgetPasswordToken", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByForgetPasswordToken indicates an expected call of FindUserByForgetPasswordToken.
func (mr *MockQueriesMockRecorder) FindUserByForgetPasswordToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByForgetPasswordToken", reflect.TypeOf((*MockQueries)(nil).FindUserByForgetPasswordToken), arg0)
}

// FindUserByID mocks base method.
func (m *MockQueries) FindUserByID(arg0 uint) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByID", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByID indicates an expected call of FindUserByID.
func (mr *MockQueriesMockRecorder) FindUserByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByID", reflect.TypeOf((*MockQueries)(nil).FindUserByID), arg0)
}

// GetLeaves mocks base method.
func (m *MockQueries) GetLeaves(arg0 *model.User) []queries.LeaveType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeaves", arg0)
	ret0, _ := ret[0].([]queries.LeaveType)
	return ret0
}

// GetLeaves indicates an expected call of GetLeaves.
func (mr *MockQueriesMockRecorder) GetLeaves(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeaves", reflect.TypeOf((*MockQueries)(nil).GetLeaves), arg0)
}

// UpdateUserCompanyID mocks base method.
func (m *MockQueries) UpdateUserCompanyID(arg0 model.User, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserCompanyID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserCompanyID indicates an expected call of UpdateUserCompanyID.
func (mr *MockQueriesMockRecorder) UpdateUserCompanyID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserCompanyID", reflect.TypeOf((*MockQueries)(nil).UpdateUserCompanyID), arg0, arg1)
}

// UpdateUserForgetPasswordToken mocks base method.
func (m *MockQueries) UpdateUserForgetPasswordToken(arg0 model.User, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserForgetPasswordToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserForgetPasswordToken indicates an expected call of UpdateUserForgetPasswordToken.
func (mr *MockQueriesMockRecorder) UpdateUserForgetPasswordToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserForgetPasswordToken", reflect.TypeOf((*MockQueries)(nil).UpdateUserForgetPasswordToken), arg0, arg1)
}

// UpdateUserPassword mocks base method.
func (m *MockQueries) UpdateUserPassword(arg0 model.User, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockQueriesMockRecorder) UpdateUserPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockQueries)(nil).UpdateUserPassword), arg0, arg1)
}
