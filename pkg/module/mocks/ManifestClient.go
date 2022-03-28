// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	manifestival "github.com/manifestival/manifestival"
	mock "github.com/stretchr/testify/mock"
)

// ManifestClient is an autogenerated mock type for the ManifestClient type
type ManifestClient struct {
	mock.Mock
}

// LoadFromSource provides a mock function with given fields: ctx, source, path
func (_m *ManifestClient) LoadFromSource(ctx context.Context, source string, path string) (manifestival.Manifest, error) {
	ret := _m.Called(ctx, source, path)

	var r0 manifestival.Manifest
	if rf, ok := ret.Get(0).(func(context.Context, string, string) manifestival.Manifest); ok {
		r0 = rf(ctx, source, path)
	} else {
		r0 = ret.Get(0).(manifestival.Manifest)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, source, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
