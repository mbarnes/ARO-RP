// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/ARO-RP/pkg/metrics (interfaces: Interface)

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// EmitFloat mocks base method.
func (m *MockInterface) EmitFloat(arg0 string, arg1 float64, arg2 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EmitFloat", arg0, arg1, arg2)
}

// EmitFloat indicates an expected call of EmitFloat.
func (mr *MockInterfaceMockRecorder) EmitFloat(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmitFloat", reflect.TypeOf((*MockInterface)(nil).EmitFloat), arg0, arg1, arg2)
}

// EmitGauge mocks base method.
func (m *MockInterface) EmitGauge(arg0 string, arg1 int64, arg2 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EmitGauge", arg0, arg1, arg2)
}

// EmitGauge indicates an expected call of EmitGauge.
func (mr *MockInterfaceMockRecorder) EmitGauge(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmitGauge", reflect.TypeOf((*MockInterface)(nil).EmitGauge), arg0, arg1, arg2)
}
