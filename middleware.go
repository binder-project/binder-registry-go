package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// logger wraps a Handler with some quick logging (after request)
func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"[Registry] %s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func contentTypeJSON(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}

// Instead of squashing a potential error in serializing, do it on startup to
// ensure we can write directly. This relies on DontPanicError still in lieu of
// creating a const string like
//   {"message":"Internal Server Error. Don't Panic. We will."}
var rawPanicMessage []byte

func init() {
	var err error
	rawPanicMessage, err = json.Marshal(DontPanicError)
	if err != nil {
		log.Panicf("Unable to initialize our panic message: %v", err)
		os.Exit(3)
	}
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(rawPanicMessage)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
