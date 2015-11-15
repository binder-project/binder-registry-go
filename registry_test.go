package registry

/**
 * Mostly helper utilities
 */

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/mock"
)

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
	equals(tb, 1, len(w.HeaderMap[key]))
	equals(tb, value, w.HeaderMap[key][0])
}

func contentTypeWasJSON(tb testing.TB, w *httptest.ResponseRecorder) {
	headerHas(tb, w, "Content-Type", "application/json; charset=UTF-8")
}

func matchedAPIError(tb testing.TB, expected APIErrorResponse, w *httptest.ResponseRecorder) {
	// Content type better be JSON
	contentTypeWasJSON(tb, w)
	// ...and the content better be a standardized error in JSON too!
	var apiError APIErrorResponse
	err := json.NewDecoder(w.Body).Decode(&apiError)
	ok(tb, err)
	equals(tb, expected.Message, apiError.Message)
}

/*

*/
type mockStore struct {
    mock.Mock
    registerCall map[string]RegisterReceiver
}

type RegisterReceiver struct {
    Input           Template
    OutputTemplate  Template
    OutputError     error
    Actual          Template
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
    receiver := _m.registerCall[tmpl.Name]
    receiver.Actual = tmpl

    return receiver.OutputTemplate, receiver.OutputError
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

func (_m *mockStore) UpdateTemplate(name string, update map[string]interface{}) (Template, error) {
    ret := _m.Called(name, update)

    var r0 Template
    if rf, ok := ret.Get(0).(func(string, map[string]interface{}) Template); ok {
        r0 = rf(name, update)
    } else {
        r0 = ret.Get(0).(Template)
    }

    var r1 error
    if rf, ok := ret.Get(1).(func(string, map[string]interface{}) error); ok {
        r1 = rf(name, update)
    } else {
        r1 = ret.Error(1)
    }

    return r0, r1
}

