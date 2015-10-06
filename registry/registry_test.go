package registry

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

/**
 * this go file only contains helper utilities for testing
 */

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
