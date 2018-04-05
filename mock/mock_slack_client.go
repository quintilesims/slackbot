// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/slackbot/utils (interfaces: SlackClient)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	slack "github.com/nlopes/slack"
)

// MockSlackClient is a mock of SlackClient interface
type MockSlackClient struct {
	ctrl     *gomock.Controller
	recorder *MockSlackClientMockRecorder
}

// MockSlackClientMockRecorder is the mock recorder for MockSlackClient
type MockSlackClientMockRecorder struct {
	mock *MockSlackClient
}

// NewMockSlackClient creates a new mock instance
func NewMockSlackClient(ctrl *gomock.Controller) *MockSlackClient {
	mock := &MockSlackClient{ctrl: ctrl}
	mock.recorder = &MockSlackClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSlackClient) EXPECT() *MockSlackClientMockRecorder {
	return m.recorder
}

// DeleteMessage mocks base method
func (m *MockSlackClient) DeleteMessage(arg0, arg1 string) (string, string, error) {
	ret := m.ctrl.Call(m, "DeleteMessage", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteMessage indicates an expected call of DeleteMessage
func (mr *MockSlackClientMockRecorder) DeleteMessage(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockSlackClient)(nil).DeleteMessage), arg0, arg1)
}

// GetChannelHistory mocks base method
func (m *MockSlackClient) GetChannelHistory(arg0 string, arg1 slack.HistoryParameters) (*slack.History, error) {
	ret := m.ctrl.Call(m, "GetChannelHistory", arg0, arg1)
	ret0, _ := ret[0].(*slack.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannelHistory indicates an expected call of GetChannelHistory
func (mr *MockSlackClientMockRecorder) GetChannelHistory(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelHistory", reflect.TypeOf((*MockSlackClient)(nil).GetChannelHistory), arg0, arg1)
}

// GetGroupHistory mocks base method
func (m *MockSlackClient) GetGroupHistory(arg0 string, arg1 slack.HistoryParameters) (*slack.History, error) {
	ret := m.ctrl.Call(m, "GetGroupHistory", arg0, arg1)
	ret0, _ := ret[0].(*slack.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupHistory indicates an expected call of GetGroupHistory
func (mr *MockSlackClientMockRecorder) GetGroupHistory(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupHistory", reflect.TypeOf((*MockSlackClient)(nil).GetGroupHistory), arg0, arg1)
}

// GetIMHistory mocks base method
func (m *MockSlackClient) GetIMHistory(arg0 string, arg1 slack.HistoryParameters) (*slack.History, error) {
	ret := m.ctrl.Call(m, "GetIMHistory", arg0, arg1)
	ret0, _ := ret[0].(*slack.History)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIMHistory indicates an expected call of GetIMHistory
func (mr *MockSlackClientMockRecorder) GetIMHistory(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIMHistory", reflect.TypeOf((*MockSlackClient)(nil).GetIMHistory), arg0, arg1)
}

// SendMessage mocks base method
func (m *MockSlackClient) SendMessage(arg0 string, arg1 ...slack.MsgOption) (string, string, string, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessage", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SendMessage indicates an expected call of SendMessage
func (mr *MockSlackClientMockRecorder) SendMessage(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockSlackClient)(nil).SendMessage), varargs...)
}
