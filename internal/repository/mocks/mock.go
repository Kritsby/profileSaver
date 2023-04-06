// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "dev/profileSaver/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(u model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", u)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), u)
}

// DeleteUser mocks base method.
func (m *MockRepository) DeleteUser(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockRepositoryMockRecorder) DeleteUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockRepository)(nil).DeleteUser), id)
}

// GetAllUsers mocks base method.
func (m *MockRepository) GetAllUsers() []model.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]model.User)
	return ret0
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockRepositoryMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockRepository)(nil).GetAllUsers))
}

// GetUserByID mocks base method.
func (m *MockRepository) GetUserByID(id string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockRepository)(nil).GetUserByID), id)
}

// GetUserByName mocks base method.
func (m *MockRepository) GetUserByName(name string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", name)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockRepositoryMockRecorder) GetUserByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockRepository)(nil).GetUserByName), name)
}

// HashPass mocks base method.
func (m *MockRepository) HashPass(password, salt []byte) ([]byte, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPass", password, salt)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// HashPass indicates an expected call of HashPass.
func (mr *MockRepositoryMockRecorder) HashPass(password, salt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPass", reflect.TypeOf((*MockRepository)(nil).HashPass), password, salt)
}

// UpdateUser mocks base method.
func (m *MockRepository) UpdateUser(u model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", u)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockRepositoryMockRecorder) UpdateUser(u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockRepository)(nil).UpdateUser), u)
}