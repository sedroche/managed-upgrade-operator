// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/managed-upgrade-operator/pkg/metrics (interfaces: Metrics)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	metrics "github.com/openshift/managed-upgrade-operator/pkg/metrics"
	reflect "reflect"
	time "time"
)

// MockMetrics is a mock of Metrics interface
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// IsAlertFiring mocks base method
func (m *MockMetrics) IsAlertFiring(arg0 string, arg1 []string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAlertFiring", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAlertFiring indicates an expected call of IsAlertFiring
func (mr *MockMetricsMockRecorder) IsAlertFiring(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAlertFiring", reflect.TypeOf((*MockMetrics)(nil).IsAlertFiring), arg0, arg1)
}

// IsMetricControlPlaneEndTimeSet mocks base method
func (m *MockMetrics) IsMetricControlPlaneEndTimeSet(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMetricControlPlaneEndTimeSet", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMetricControlPlaneEndTimeSet indicates an expected call of IsMetricControlPlaneEndTimeSet
func (mr *MockMetricsMockRecorder) IsMetricControlPlaneEndTimeSet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMetricControlPlaneEndTimeSet", reflect.TypeOf((*MockMetrics)(nil).IsMetricControlPlaneEndTimeSet), arg0, arg1)
}

// IsMetricNodeUpgradeEndTimeSet mocks base method
func (m *MockMetrics) IsMetricNodeUpgradeEndTimeSet(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMetricNodeUpgradeEndTimeSet", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMetricNodeUpgradeEndTimeSet indicates an expected call of IsMetricNodeUpgradeEndTimeSet
func (mr *MockMetricsMockRecorder) IsMetricNodeUpgradeEndTimeSet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMetricNodeUpgradeEndTimeSet", reflect.TypeOf((*MockMetrics)(nil).IsMetricNodeUpgradeEndTimeSet), arg0, arg1)
}

// IsMetricUpgradeStartTimeSet mocks base method
func (m *MockMetrics) IsMetricUpgradeStartTimeSet(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMetricUpgradeStartTimeSet", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMetricUpgradeStartTimeSet indicates an expected call of IsMetricUpgradeStartTimeSet
func (mr *MockMetricsMockRecorder) IsMetricUpgradeStartTimeSet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMetricUpgradeStartTimeSet", reflect.TypeOf((*MockMetrics)(nil).IsMetricUpgradeStartTimeSet), arg0, arg1)
}

// Query mocks base method
func (m *MockMetrics) Query(arg0 string) (*metrics.AlertResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(*metrics.AlertResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockMetricsMockRecorder) Query(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockMetrics)(nil).Query), arg0)
}

// ResetAllMetricNodeDrainFailed mocks base method
func (m *MockMetrics) ResetAllMetricNodeDrainFailed() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetAllMetricNodeDrainFailed")
}

// ResetAllMetricNodeDrainFailed indicates an expected call of ResetAllMetricNodeDrainFailed
func (mr *MockMetricsMockRecorder) ResetAllMetricNodeDrainFailed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetAllMetricNodeDrainFailed", reflect.TypeOf((*MockMetrics)(nil).ResetAllMetricNodeDrainFailed))
}

// ResetMetricNodeDrainFailed mocks base method
func (m *MockMetrics) ResetMetricNodeDrainFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetMetricNodeDrainFailed", arg0)
}

// ResetMetricNodeDrainFailed indicates an expected call of ResetMetricNodeDrainFailed
func (mr *MockMetricsMockRecorder) ResetMetricNodeDrainFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetMetricNodeDrainFailed", reflect.TypeOf((*MockMetrics)(nil).ResetMetricNodeDrainFailed), arg0)
}

// ResetMetricUpgradeControlPlaneTimeout mocks base method
func (m *MockMetrics) ResetMetricUpgradeControlPlaneTimeout(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetMetricUpgradeControlPlaneTimeout", arg0, arg1)
}

// ResetMetricUpgradeControlPlaneTimeout indicates an expected call of ResetMetricUpgradeControlPlaneTimeout
func (mr *MockMetricsMockRecorder) ResetMetricUpgradeControlPlaneTimeout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetMetricUpgradeControlPlaneTimeout", reflect.TypeOf((*MockMetrics)(nil).ResetMetricUpgradeControlPlaneTimeout), arg0, arg1)
}

// ResetMetricUpgradeWorkerTimeout mocks base method
func (m *MockMetrics) ResetMetricUpgradeWorkerTimeout(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetMetricUpgradeWorkerTimeout", arg0, arg1)
}

// ResetMetricUpgradeWorkerTimeout indicates an expected call of ResetMetricUpgradeWorkerTimeout
func (mr *MockMetricsMockRecorder) ResetMetricUpgradeWorkerTimeout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetMetricUpgradeWorkerTimeout", reflect.TypeOf((*MockMetrics)(nil).ResetMetricUpgradeWorkerTimeout), arg0, arg1)
}

// ResetMetrics mocks base method
func (m *MockMetrics) ResetMetrics() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetMetrics")
}

// ResetMetrics indicates an expected call of ResetMetrics
func (mr *MockMetricsMockRecorder) ResetMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetMetrics", reflect.TypeOf((*MockMetrics)(nil).ResetMetrics))
}

// UpdateMetricClusterCheckFailed mocks base method
func (m *MockMetrics) UpdateMetricClusterCheckFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricClusterCheckFailed", arg0)
}

// UpdateMetricClusterCheckFailed indicates an expected call of UpdateMetricClusterCheckFailed
func (mr *MockMetricsMockRecorder) UpdateMetricClusterCheckFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricClusterCheckFailed", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricClusterCheckFailed), arg0)
}

// UpdateMetricClusterCheckSucceeded mocks base method
func (m *MockMetrics) UpdateMetricClusterCheckSucceeded(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricClusterCheckSucceeded", arg0)
}

// UpdateMetricClusterCheckSucceeded indicates an expected call of UpdateMetricClusterCheckSucceeded
func (mr *MockMetricsMockRecorder) UpdateMetricClusterCheckSucceeded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricClusterCheckSucceeded", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricClusterCheckSucceeded), arg0)
}

// UpdateMetricClusterVerificationFailed mocks base method
func (m *MockMetrics) UpdateMetricClusterVerificationFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricClusterVerificationFailed", arg0)
}

// UpdateMetricClusterVerificationFailed indicates an expected call of UpdateMetricClusterVerificationFailed
func (mr *MockMetricsMockRecorder) UpdateMetricClusterVerificationFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricClusterVerificationFailed", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricClusterVerificationFailed), arg0)
}

// UpdateMetricClusterVerificationSucceeded mocks base method
func (m *MockMetrics) UpdateMetricClusterVerificationSucceeded(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricClusterVerificationSucceeded", arg0)
}

// UpdateMetricClusterVerificationSucceeded indicates an expected call of UpdateMetricClusterVerificationSucceeded
func (mr *MockMetricsMockRecorder) UpdateMetricClusterVerificationSucceeded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricClusterVerificationSucceeded", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricClusterVerificationSucceeded), arg0)
}

// UpdateMetricControlPlaneEndTime mocks base method
func (m *MockMetrics) UpdateMetricControlPlaneEndTime(arg0 time.Time, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricControlPlaneEndTime", arg0, arg1, arg2)
}

// UpdateMetricControlPlaneEndTime indicates an expected call of UpdateMetricControlPlaneEndTime
func (mr *MockMetricsMockRecorder) UpdateMetricControlPlaneEndTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricControlPlaneEndTime", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricControlPlaneEndTime), arg0, arg1, arg2)
}

// UpdateMetricNodeDrainFailed mocks base method
func (m *MockMetrics) UpdateMetricNodeDrainFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricNodeDrainFailed", arg0)
}

// UpdateMetricNodeDrainFailed indicates an expected call of UpdateMetricNodeDrainFailed
func (mr *MockMetricsMockRecorder) UpdateMetricNodeDrainFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricNodeDrainFailed", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricNodeDrainFailed), arg0)
}

// UpdateMetricNodeUpgradeEndTime mocks base method
func (m *MockMetrics) UpdateMetricNodeUpgradeEndTime(arg0 time.Time, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricNodeUpgradeEndTime", arg0, arg1, arg2)
}

// UpdateMetricNodeUpgradeEndTime indicates an expected call of UpdateMetricNodeUpgradeEndTime
func (mr *MockMetricsMockRecorder) UpdateMetricNodeUpgradeEndTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricNodeUpgradeEndTime", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricNodeUpgradeEndTime), arg0, arg1, arg2)
}

// UpdateMetricScalingFailed mocks base method
func (m *MockMetrics) UpdateMetricScalingFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricScalingFailed", arg0)
}

// UpdateMetricScalingFailed indicates an expected call of UpdateMetricScalingFailed
func (mr *MockMetricsMockRecorder) UpdateMetricScalingFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricScalingFailed", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricScalingFailed), arg0)
}

// UpdateMetricScalingSucceeded mocks base method
func (m *MockMetrics) UpdateMetricScalingSucceeded(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricScalingSucceeded", arg0)
}

// UpdateMetricScalingSucceeded indicates an expected call of UpdateMetricScalingSucceeded
func (mr *MockMetricsMockRecorder) UpdateMetricScalingSucceeded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricScalingSucceeded", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricScalingSucceeded), arg0)
}

// UpdateMetricUpgradeControlPlaneTimeout mocks base method
func (m *MockMetrics) UpdateMetricUpgradeControlPlaneTimeout(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricUpgradeControlPlaneTimeout", arg0, arg1)
}

// UpdateMetricUpgradeControlPlaneTimeout indicates an expected call of UpdateMetricUpgradeControlPlaneTimeout
func (mr *MockMetricsMockRecorder) UpdateMetricUpgradeControlPlaneTimeout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricUpgradeControlPlaneTimeout", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricUpgradeControlPlaneTimeout), arg0, arg1)
}

// UpdateMetricUpgradeStartTime mocks base method
func (m *MockMetrics) UpdateMetricUpgradeStartTime(arg0 time.Time, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricUpgradeStartTime", arg0, arg1, arg2)
}

// UpdateMetricUpgradeStartTime indicates an expected call of UpdateMetricUpgradeStartTime
func (mr *MockMetricsMockRecorder) UpdateMetricUpgradeStartTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricUpgradeStartTime", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricUpgradeStartTime), arg0, arg1, arg2)
}

// UpdateMetricUpgradeWindowBreached mocks base method
func (m *MockMetrics) UpdateMetricUpgradeWindowBreached(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricUpgradeWindowBreached", arg0)
}

// UpdateMetricUpgradeWindowBreached indicates an expected call of UpdateMetricUpgradeWindowBreached
func (mr *MockMetricsMockRecorder) UpdateMetricUpgradeWindowBreached(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricUpgradeWindowBreached", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricUpgradeWindowBreached), arg0)
}

// UpdateMetricUpgradeWindowNotBreached mocks base method
func (m *MockMetrics) UpdateMetricUpgradeWindowNotBreached(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricUpgradeWindowNotBreached", arg0)
}

// UpdateMetricUpgradeWindowNotBreached indicates an expected call of UpdateMetricUpgradeWindowNotBreached
func (mr *MockMetricsMockRecorder) UpdateMetricUpgradeWindowNotBreached(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricUpgradeWindowNotBreached", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricUpgradeWindowNotBreached), arg0)
}

// UpdateMetricUpgradeWorkerTimeout mocks base method
func (m *MockMetrics) UpdateMetricUpgradeWorkerTimeout(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricUpgradeWorkerTimeout", arg0, arg1)
}

// UpdateMetricUpgradeWorkerTimeout indicates an expected call of UpdateMetricUpgradeWorkerTimeout
func (mr *MockMetricsMockRecorder) UpdateMetricUpgradeWorkerTimeout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricUpgradeWorkerTimeout", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricUpgradeWorkerTimeout), arg0, arg1)
}

// UpdateMetricValidationFailed mocks base method
func (m *MockMetrics) UpdateMetricValidationFailed(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricValidationFailed", arg0)
}

// UpdateMetricValidationFailed indicates an expected call of UpdateMetricValidationFailed
func (mr *MockMetricsMockRecorder) UpdateMetricValidationFailed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricValidationFailed", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricValidationFailed), arg0)
}

// UpdateMetricValidationSucceeded mocks base method
func (m *MockMetrics) UpdateMetricValidationSucceeded(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateMetricValidationSucceeded", arg0)
}

// UpdateMetricValidationSucceeded indicates an expected call of UpdateMetricValidationSucceeded
func (mr *MockMetricsMockRecorder) UpdateMetricValidationSucceeded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricValidationSucceeded", reflect.TypeOf((*MockMetrics)(nil).UpdateMetricValidationSucceeded), arg0)
}
