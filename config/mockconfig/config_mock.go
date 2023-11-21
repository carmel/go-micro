//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package mockconfig is a generated GoMock package.
package mockconfig

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "go-micro/config"
)

// MockUnmarshaler is a mock of Unmarshaler interface
type MockUnmarshaler struct {
	ctrl     *gomock.Controller
	recorder *MockUnmarshalerMockRecorder
}

// MockUnmarshalerMockRecorder is the mock recorder for MockUnmarshaler
type MockUnmarshalerMockRecorder struct {
	mock *MockUnmarshaler
}

// NewMockUnmarshaler creates a new mock instance
func NewMockUnmarshaler(ctrl *gomock.Controller) *MockUnmarshaler {
	mock := &MockUnmarshaler{ctrl: ctrl}
	mock.recorder = &MockUnmarshalerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUnmarshaler) EXPECT() *MockUnmarshalerMockRecorder {
	return m.recorder
}

// Unmarshal mocks base method
func (m *MockUnmarshaler) Unmarshal(data []byte, value interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", data, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal
func (mr *MockUnmarshalerMockRecorder) Unmarshal(data, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockUnmarshaler)(nil).Unmarshal), data, value)
}

// MockKVConfig is a mock of KVConfig interface
type MockKVConfig struct {
	ctrl     *gomock.Controller
	recorder *MockKVConfigMockRecorder
}

// MockKVConfigMockRecorder is the mock recorder for MockKVConfig
type MockKVConfigMockRecorder struct {
	mock *MockKVConfig
}

// NewMockKVConfig creates a new mock instance
func NewMockKVConfig(ctrl *gomock.Controller) *MockKVConfig {
	mock := &MockKVConfig{ctrl: ctrl}
	mock.recorder = &MockKVConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKVConfig) EXPECT() *MockKVConfigMockRecorder {
	return m.recorder
}

// Put mocks base method
func (m *MockKVConfig) Put(ctx context.Context, key, val string, opts ...config.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key, val}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Put", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockKVConfigMockRecorder) Put(ctx, key, val interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, val}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockKVConfig)(nil).Put), varargs...)
}

// Get mocks base method
func (m *MockKVConfig) Get(ctx context.Context, key string, opts ...config.Option) (config.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(config.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKVConfigMockRecorder) Get(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKVConfig)(nil).Get), varargs...)
}

// Del mocks base method
func (m *MockKVConfig) Del(ctx context.Context, key string, opts ...config.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Del", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del
func (mr *MockKVConfigMockRecorder) Del(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockKVConfig)(nil).Del), varargs...)
}

// Watch mocks base method
func (m *MockKVConfig) Watch(ctx context.Context, key string, opts ...config.Option) (<-chan config.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(<-chan config.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockKVConfigMockRecorder) Watch(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockKVConfig)(nil).Watch), varargs...)
}

// Name mocks base method
func (m *MockKVConfig) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockKVConfigMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockKVConfig)(nil).Name))
}

// MockResponse is a mock of Response interface
type MockResponse struct {
	ctrl     *gomock.Controller
	recorder *MockResponseMockRecorder
}

// MockResponseMockRecorder is the mock recorder for MockResponse
type MockResponseMockRecorder struct {
	mock *MockResponse
}

// NewMockResponse creates a new mock instance
func NewMockResponse(ctrl *gomock.Controller) *MockResponse {
	mock := &MockResponse{ctrl: ctrl}
	mock.recorder = &MockResponseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResponse) EXPECT() *MockResponseMockRecorder {
	return m.recorder
}

// Value mocks base method
func (m *MockResponse) Value() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value")
	ret0, _ := ret[0].(string)
	return ret0
}

// Value indicates an expected call of Value
func (mr *MockResponseMockRecorder) Value() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockResponse)(nil).Value))
}

// MetaData mocks base method
func (m *MockResponse) MetaData() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MetaData")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// MetaData indicates an expected call of MetaData
func (mr *MockResponseMockRecorder) MetaData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MetaData", reflect.TypeOf((*MockResponse)(nil).MetaData))
}

// Event mocks base method
func (m *MockResponse) Event() config.EventType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Event")
	ret0, _ := ret[0].(config.EventType)
	return ret0
}

// Event indicates an expected call of Event
func (mr *MockResponseMockRecorder) Event() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Event", reflect.TypeOf((*MockResponse)(nil).Event))
}

// MockKV is a mock of KV interface
type MockKV struct {
	ctrl     *gomock.Controller
	recorder *MockKVMockRecorder
}

// MockKVMockRecorder is the mock recorder for MockKV
type MockKVMockRecorder struct {
	mock *MockKV
}

// NewMockKV creates a new mock instance
func NewMockKV(ctrl *gomock.Controller) *MockKV {
	mock := &MockKV{ctrl: ctrl}
	mock.recorder = &MockKVMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKV) EXPECT() *MockKVMockRecorder {
	return m.recorder
}

// Put mocks base method
func (m *MockKV) Put(ctx context.Context, key, val string, opts ...config.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key, val}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Put", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockKVMockRecorder) Put(ctx, key, val interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, val}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockKV)(nil).Put), varargs...)
}

// Get mocks base method
func (m *MockKV) Get(ctx context.Context, key string, opts ...config.Option) (config.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(config.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKVMockRecorder) Get(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKV)(nil).Get), varargs...)
}

// Del mocks base method
func (m *MockKV) Del(ctx context.Context, key string, opts ...config.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Del", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del
func (mr *MockKVMockRecorder) Del(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockKV)(nil).Del), varargs...)
}

// MockWatcher is a mock of Watcher interface
type MockWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher
type MockWatcherMockRecorder struct {
	mock *MockWatcher
}

// NewMockWatcher creates a new mock instance
func NewMockWatcher(ctrl *gomock.Controller) *MockWatcher {
	mock := &MockWatcher{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWatcher) EXPECT() *MockWatcherMockRecorder {
	return m.recorder
}

// Watch mocks base method
func (m *MockWatcher) Watch(ctx context.Context, key string, opts ...config.Option) (<-chan config.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(<-chan config.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockWatcherMockRecorder) Watch(ctx, key interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockWatcher)(nil).Watch), varargs...)
}

// MockLoader is a mock of Loader interface
type MockLoader struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderMockRecorder
}

// MockLoaderMockRecorder is the mock recorder for MockLoader
type MockLoaderMockRecorder struct {
	mock *MockLoader
}

// NewMockLoader creates a new mock instance
func NewMockLoader(ctrl *gomock.Controller) *MockLoader {
	mock := &MockLoader{ctrl: ctrl}
	mock.recorder = &MockLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoader) EXPECT() *MockLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockLoader) Load(arg0 string) (config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].(config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load
func (mr *MockLoaderMockRecorder) Load(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockLoader)(nil).Load), arg0)
}

// Reload mocks base method
func (m *MockLoader) Reload(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reload indicates an expected call of Reload
func (mr *MockLoaderMockRecorder) Reload(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*MockLoader)(nil).Reload), arg0)
}

// MockConfig is a mock of Config interface
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockConfig) Load() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load")
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load
func (mr *MockConfigMockRecorder) Load() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockConfig)(nil).Load))
}

// Reload mocks base method
func (m *MockConfig) Reload() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reload")
}

// Reload indicates an expected call of Reload
func (mr *MockConfigMockRecorder) Reload() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*MockConfig)(nil).Reload))
}

// Get mocks base method
func (m *MockConfig) Get(arg0 string, arg1 interface{}) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockConfigMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfig)(nil).Get), arg0, arg1)
}

// Unmarshal mocks base method
func (m *MockConfig) Unmarshal(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal
func (mr *MockConfigMockRecorder) Unmarshal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockConfig)(nil).Unmarshal), arg0)
}

// IsSet mocks base method
func (m *MockConfig) IsSet(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSet", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSet indicates an expected call of IsSet
func (mr *MockConfigMockRecorder) IsSet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSet", reflect.TypeOf((*MockConfig)(nil).IsSet), arg0)
}

// GetInt mocks base method
func (m *MockConfig) GetInt(arg0 string, arg1 int) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", arg0, arg1)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetInt indicates an expected call of GetInt
func (mr *MockConfigMockRecorder) GetInt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockConfig)(nil).GetInt), arg0, arg1)
}

// GetInt32 mocks base method
func (m *MockConfig) GetInt32(arg0 string, arg1 int32) int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", arg0, arg1)
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetInt32 indicates an expected call of GetInt32
func (mr *MockConfigMockRecorder) GetInt32(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockConfig)(nil).GetInt32), arg0, arg1)
}

// GetInt64 mocks base method
func (m *MockConfig) GetInt64(arg0 string, arg1 int64) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt64", arg0, arg1)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetInt64 indicates an expected call of GetInt64
func (mr *MockConfigMockRecorder) GetInt64(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt64", reflect.TypeOf((*MockConfig)(nil).GetInt64), arg0, arg1)
}

// GetUint mocks base method
func (m *MockConfig) GetUint(arg0 string, arg1 uint) uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint", arg0, arg1)
	ret0, _ := ret[0].(uint)
	return ret0
}

// GetUint indicates an expected call of GetUint
func (mr *MockConfigMockRecorder) GetUint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint", reflect.TypeOf((*MockConfig)(nil).GetUint), arg0, arg1)
}

// GetUint32 mocks base method
func (m *MockConfig) GetUint32(arg0 string, arg1 uint32) uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint32", arg0, arg1)
	ret0, _ := ret[0].(uint32)
	return ret0
}

// GetUint32 indicates an expected call of GetUint32
func (mr *MockConfigMockRecorder) GetUint32(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint32", reflect.TypeOf((*MockConfig)(nil).GetUint32), arg0, arg1)
}

// GetUint64 mocks base method
func (m *MockConfig) GetUint64(arg0 string, arg1 uint64) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint64", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetUint64 indicates an expected call of GetUint64
func (mr *MockConfigMockRecorder) GetUint64(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint64", reflect.TypeOf((*MockConfig)(nil).GetUint64), arg0, arg1)
}

// GetFloat32 mocks base method
func (m *MockConfig) GetFloat32(arg0 string, arg1 float32) float32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloat32", arg0, arg1)
	ret0, _ := ret[0].(float32)
	return ret0
}

// GetFloat32 indicates an expected call of GetFloat32
func (mr *MockConfigMockRecorder) GetFloat32(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloat32", reflect.TypeOf((*MockConfig)(nil).GetFloat32), arg0, arg1)
}

// GetFloat64 mocks base method
func (m *MockConfig) GetFloat64(arg0 string, arg1 float64) float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloat64", arg0, arg1)
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetFloat64 indicates an expected call of GetFloat64
func (mr *MockConfigMockRecorder) GetFloat64(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloat64", reflect.TypeOf((*MockConfig)(nil).GetFloat64), arg0, arg1)
}

// GetString mocks base method
func (m *MockConfig) GetString(arg0, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString
func (mr *MockConfigMockRecorder) GetString(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockConfig)(nil).GetString), arg0, arg1)
}

// GetBool mocks base method
func (m *MockConfig) GetBool(arg0 string, arg1 bool) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBool", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetBool indicates an expected call of GetBool
func (mr *MockConfigMockRecorder) GetBool(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBool", reflect.TypeOf((*MockConfig)(nil).GetBool), arg0, arg1)
}

// Bytes mocks base method
func (m *MockConfig) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes
func (mr *MockConfigMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockConfig)(nil).Bytes))
}

// MockDataProvider is a mock of DataProvider interface
type MockDataProvider struct {
	ctrl     *gomock.Controller
	recorder *MockDataProviderMockRecorder
}

// MockDataProviderMockRecorder is the mock recorder for MockDataProvider
type MockDataProviderMockRecorder struct {
	mock *MockDataProvider
}

// NewMockDataProvider creates a new mock instance
func NewMockDataProvider(ctrl *gomock.Controller) *MockDataProvider {
	mock := &MockDataProvider{ctrl: ctrl}
	mock.recorder = &MockDataProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataProvider) EXPECT() *MockDataProviderMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockDataProvider) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockDataProviderMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockDataProvider)(nil).Name))
}

// Read mocks base method
func (m *MockDataProvider) Read(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockDataProviderMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockDataProvider)(nil).Read), arg0)
}

// Watch mocks base method
func (m *MockDataProvider) Watch(arg0 config.ProviderCallback) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0)
}

// Watch indicates an expected call of Watch
func (mr *MockDataProviderMockRecorder) Watch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockDataProvider)(nil).Watch), arg0)
}

// MockCodec is a mock of Codec interface
type MockCodec struct {
	ctrl     *gomock.Controller
	recorder *MockCodecMockRecorder
}

// MockCodecMockRecorder is the mock recorder for MockCodec
type MockCodecMockRecorder struct {
	mock *MockCodec
}

// NewMockCodec creates a new mock instance
func NewMockCodec(ctrl *gomock.Controller) *MockCodec {
	mock := &MockCodec{ctrl: ctrl}
	mock.recorder = &MockCodecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCodec) EXPECT() *MockCodecMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockCodec) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockCodecMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockCodec)(nil).Name))
}

// Unmarshal mocks base method
func (m *MockCodec) Unmarshal(arg0 []byte, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal
func (mr *MockCodecMockRecorder) Unmarshal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockCodec)(nil).Unmarshal), arg0, arg1)
}
