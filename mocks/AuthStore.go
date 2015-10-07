package mocks

import "github.com/stretchr/testify/mock"

import "net/http"

type AuthStore struct {
	mock.Mock
}

func (_m *AuthStore) Authorize(inner http.Handler) http.Handler {
	ret := _m.Called(inner)

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func(http.Handler) http.Handler); ok {
		r0 = rf(inner)
	} else {
		r0 = ret.Get(0).(http.Handler)
	}

	return r0
}
