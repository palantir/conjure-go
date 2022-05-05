// Copyright (c) 2022 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by godel-mockgen-plugin. DO NOT EDIT.
// Configure in "godel/config/mockgen-plugin.yml" and regenerate with "./godelw mockgen".

package api_mock

import (
	context "context"
	io "io"

	api "github.com/palantir/conjure-go/v6/integration_test/testgenerated/cli/api"
	bearertoken "github.com/palantir/pkg/bearertoken"
	datetime "github.com/palantir/pkg/datetime"
	rid "github.com/palantir/pkg/rid"
	safelong "github.com/palantir/pkg/safelong"
	uuid "github.com/palantir/pkg/uuid"
	mock "github.com/stretchr/testify/mock"
)

// TestService is an autogenerated mock type for the TestService type
type TestService struct {
	mock.Mock
}

// Chan provides a mock function
func (_m *TestService) Chan(ctx context.Context, varArg string, importArg map[string]string, typeArg string, returnArg safelong.SafeLong, httpArg string, jsonArg string, reqArg string, rwArg string) error {
	ret := _m.Called(ctx, varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string, string, safelong.SafeLong, string, string, string, string) error); ok {
		r0 = rf(ctx, varArg, importArg, typeArg, returnArg, httpArg, jsonArg, reqArg, rwArg)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Echo provides a mock function
func (_m *TestService) Echo(ctx context.Context, cookieToken bearertoken.Token) error {
	ret := _m.Called(ctx, cookieToken)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bearertoken.Token) error); ok {
		r0 = rf(ctx, cookieToken)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// EchoCustomObject provides a mock function
func (_m *TestService) EchoCustomObject(ctx context.Context, bodyArg *api.CustomObject) (*api.CustomObject, error) {
	ret := _m.Called(ctx, bodyArg)
	var r0 *api.CustomObject
	if rf, ok := ret.Get(0).(func(context.Context, *api.CustomObject) *api.CustomObject); ok {
		r0 = rf(ctx, bodyArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(*api.CustomObject)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.CustomObject) error); ok {
		r1 = rf(ctx, bodyArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// EchoOptionalAlias provides a mock function
func (_m *TestService) EchoOptionalAlias(ctx context.Context, bodyArg api.OptionalIntegerAlias) (api.OptionalIntegerAlias, error) {
	ret := _m.Called(ctx, bodyArg)
	var r0 api.OptionalIntegerAlias
	if rf, ok := ret.Get(0).(func(context.Context, api.OptionalIntegerAlias) api.OptionalIntegerAlias); ok {
		r0 = rf(ctx, bodyArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(api.OptionalIntegerAlias)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, api.OptionalIntegerAlias) error); ok {
		r1 = rf(ctx, bodyArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// EchoOptionalListAlias provides a mock function
func (_m *TestService) EchoOptionalListAlias(ctx context.Context, bodyArg api.OptionalListAlias) (api.OptionalListAlias, error) {
	ret := _m.Called(ctx, bodyArg)
	var r0 api.OptionalListAlias
	if rf, ok := ret.Get(0).(func(context.Context, api.OptionalListAlias) api.OptionalListAlias); ok {
		r0 = rf(ctx, bodyArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(api.OptionalListAlias)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, api.OptionalListAlias) error); ok {
		r1 = rf(ctx, bodyArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// EchoStrings provides a mock function
func (_m *TestService) EchoStrings(ctx context.Context, bodyArg []string) ([]string, error) {
	ret := _m.Called(ctx, bodyArg)
	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(ctx, bodyArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.([]string)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, bodyArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetBinary provides a mock function
func (_m *TestService) GetBinary(ctx context.Context) (io.ReadCloser, error) {
	ret := _m.Called(ctx)
	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context) io.ReadCloser); ok {
		r0 = rf(ctx)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(io.ReadCloser)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetDateTime provides a mock function
func (_m *TestService) GetDateTime(ctx context.Context, myParamArg datetime.DateTime) (datetime.DateTime, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 datetime.DateTime
	if rf, ok := ret.Get(0).(func(context.Context, datetime.DateTime) datetime.DateTime); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(datetime.DateTime)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, datetime.DateTime) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetDouble provides a mock function
func (_m *TestService) GetDouble(ctx context.Context, myParamArg float64) (float64, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, float64) float64); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(float64)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, float64) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetListBoolean provides a mock function
func (_m *TestService) GetListBoolean(ctx context.Context, myQueryParam1Arg []bool) ([]bool, error) {
	ret := _m.Called(ctx, myQueryParam1Arg)
	var r0 []bool
	if rf, ok := ret.Get(0).(func(context.Context, []bool) []bool); ok {
		r0 = rf(ctx, myQueryParam1Arg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.([]bool)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []bool) error); ok {
		r1 = rf(ctx, myQueryParam1Arg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetOptionalBinary provides a mock function
func (_m *TestService) GetOptionalBinary(ctx context.Context) (*io.ReadCloser, error) {
	ret := _m.Called(ctx)
	var r0 *io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context) *io.ReadCloser); ok {
		r0 = rf(ctx)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(*io.ReadCloser)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetPathParam provides a mock function
func (_m *TestService) GetPathParam(ctx context.Context, authHeader bearertoken.Token, myPathParamArg string) error {
	ret := _m.Called(ctx, authHeader, myPathParamArg)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bearertoken.Token, string) error); ok {
		r0 = rf(ctx, authHeader, myPathParamArg)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// GetReserved provides a mock function
func (_m *TestService) GetReserved(ctx context.Context, confArg string, bearertokenArg string) error {
	ret := _m.Called(ctx, confArg, bearertokenArg)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, confArg, bearertokenArg)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// GetRid provides a mock function
func (_m *TestService) GetRid(ctx context.Context, myParamArg rid.ResourceIdentifier) (rid.ResourceIdentifier, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 rid.ResourceIdentifier
	if rf, ok := ret.Get(0).(func(context.Context, rid.ResourceIdentifier) rid.ResourceIdentifier); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(rid.ResourceIdentifier)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, rid.ResourceIdentifier) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetSafeLong provides a mock function
func (_m *TestService) GetSafeLong(ctx context.Context, myParamArg safelong.SafeLong) (safelong.SafeLong, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 safelong.SafeLong
	if rf, ok := ret.Get(0).(func(context.Context, safelong.SafeLong) safelong.SafeLong); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(safelong.SafeLong)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, safelong.SafeLong) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// GetUuid provides a mock function
func (_m *TestService) GetUuid(ctx context.Context, myParamArg uuid.UUID) (uuid.UUID, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) uuid.UUID); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(uuid.UUID)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PutMapStringAny provides a mock function
func (_m *TestService) PutMapStringAny(ctx context.Context, myParamArg map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(map[string]interface{})
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PutMapStringString provides a mock function
func (_m *TestService) PutMapStringString(ctx context.Context, myParamArg map[string]string) (map[string]string, error) {
	ret := _m.Called(ctx, myParamArg)
	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) map[string]string); ok {
		r0 = rf(ctx, myParamArg)
	} else if v := ret.Get(0); v != nil {
		r0 = v.(map[string]string)
	}
	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]string) error); ok {
		r1 = rf(ctx, myParamArg)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}