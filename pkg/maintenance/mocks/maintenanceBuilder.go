// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/managed-upgrade-operator/pkg/maintenance (interfaces: MaintenanceBuilder)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	maintenance "github.com/openshift/managed-upgrade-operator/pkg/maintenance"
	reflect "reflect"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockMaintenanceBuilder is a mock of MaintenanceBuilder interface
type MockMaintenanceBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockMaintenanceBuilderMockRecorder
}

// MockMaintenanceBuilderMockRecorder is the mock recorder for MockMaintenanceBuilder
type MockMaintenanceBuilderMockRecorder struct {
	mock *MockMaintenanceBuilder
}

// NewMockMaintenanceBuilder creates a new mock instance
func NewMockMaintenanceBuilder(ctrl *gomock.Controller) *MockMaintenanceBuilder {
	mock := &MockMaintenanceBuilder{ctrl: ctrl}
	mock.recorder = &MockMaintenanceBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMaintenanceBuilder) EXPECT() *MockMaintenanceBuilderMockRecorder {
	return m.recorder
}

// NewClient mocks base method
func (m *MockMaintenanceBuilder) NewClient(arg0 client.Client, arg1 maintenance.MaintenanceConfig) (maintenance.Maintenance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClient", arg0, arg1)
	ret0, _ := ret[0].(maintenance.Maintenance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClient indicates an expected call of NewClient
func (mr *MockMaintenanceBuilderMockRecorder) NewClient(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClient", reflect.TypeOf((*MockMaintenanceBuilder)(nil).NewClient), arg0, arg1)
}
