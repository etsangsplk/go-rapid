// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import remoting "github.com/casualjim/go-rapid/remoting"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Client) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Do provides a mock function with given fields: ctx, target, in
func (_m *Client) Do(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	ret := _m.Called(ctx, target, in)

	var r0 *remoting.RapidResponse
	if rf, ok := ret.Get(0).(func(context.Context, *remoting.Endpoint, *remoting.RapidRequest) *remoting.RapidResponse); ok {
		r0 = rf(ctx, target, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*remoting.RapidResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *remoting.Endpoint, *remoting.RapidRequest) error); ok {
		r1 = rf(ctx, target, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DoBestEffort provides a mock function with given fields: ctx, target, in
func (_m *Client) DoBestEffort(ctx context.Context, target *remoting.Endpoint, in *remoting.RapidRequest) (*remoting.RapidResponse, error) {
	ret := _m.Called(ctx, target, in)

	var r0 *remoting.RapidResponse
	if rf, ok := ret.Get(0).(func(context.Context, *remoting.Endpoint, *remoting.RapidRequest) *remoting.RapidResponse); ok {
		r0 = rf(ctx, target, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*remoting.RapidResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *remoting.Endpoint, *remoting.RapidRequest) error); ok {
		r1 = rf(ctx, target, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
