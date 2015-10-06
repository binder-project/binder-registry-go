package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/mock"
)

/**
 * this go file only contains helper utilities for testing
 */

var authStore TokenAuthStore
var store InMemoryStore
var registry Registry
var authorize http.Handler
var req *http.Request
var w *httptest.ResponseRecorder

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// nequals fails the test if exp is equal to act.
func nequals(tb testing.TB, exp, act interface{}) {
	if reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:expected unequal:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func headerHas(tb testing.TB, w *httptest.ResponseRecorder, key string, value string) {
	fmt.Println(w.HeaderMap[key])
	equals(tb, 1, len(w.HeaderMap[key]))
	equals(tb, value, w.HeaderMap[key][0])
}

func matchedAPIError(tb testing.TB, expected APIErrorResponse, w *httptest.ResponseRecorder) {
	// Content type better be JSON
	headerHas(tb, w, "Content-Type", "application/json; charset=UTF-8")
	// ...and the content better be a standardized error in JSON too!
	var apiError APIErrorResponse
	err := json.NewDecoder(w.Body).Decode(&apiError)
	ok(tb, err)
	equals(tb, expected.Message, apiError.Message)
}

/**
autogenerated with mockery -name=store
**/
type mockStore struct {
	mock.Mock
}

func (_m *mockStore) GetTemplate(name string) (Template, error) {
	ret := _m.Called(name)

	var r0 Template
	if rf, ok := ret.Get(0).(func(string) Template); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *mockStore) RegisterTemplate(tmpl Template) (Template, error) {
	ret := _m.Called(tmpl)

	var r0 Template
	if rf, ok := ret.Get(0).(func(Template) Template); ok {
		r0 = rf(tmpl)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Template) error); ok {
		r1 = rf(tmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *mockStore) ListTemplates() ([]Template, error) {
	ret := _m.Called()

	var r0 []Template
	if rf, ok := ret.Get(0).(func() []Template); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *mockStore) UpdateTemplate(tmpl Template) (Template, error) {
	ret := _m.Called(tmpl)

	var r0 Template
	if rf, ok := ret.Get(0).(func(Template) Template); ok {
		r0 = rf(tmpl)
	} else {
		r0 = ret.Get(0).(Template)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Template) error); ok {
		r1 = rf(tmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
