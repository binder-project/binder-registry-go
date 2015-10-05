package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// Authenticated probes for an Authorization header
func (ctxt RegistryContext) Authenticated(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, hasAuthHeader := r.Header["Authorization"]

		if !hasAuthHeader || len(authHeader) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			err := APIErrorResponse{Message: "Authorization header not set. Should be of format 'Authorization: token key'"}
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		// QUESTION: What happens when there are multiple Authorization headers?
		tokenLine := authHeader[0]
		authHeader = strings.SplitN(tokenLine, " ", 3)

		if len(authHeader) != 2 || strings.ToLower(authHeader[0]) != "token" {
			w.WriteHeader(http.StatusUnauthorized)
			err := APIErrorResponse{Message: "Token field not present in Authorization header. Should be of format 'Authorization: token key'"}
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		putativeToken := authHeader[1]

		// TODO: HMAC SHA between shared secret + values, minding timing attacks
		// TODO #2: See if go has some of this built in nicely

		if putativeToken != ctxt.Token {
			w.WriteHeader(http.StatusUnauthorized)
			err := APIErrorResponse{Message: "Invalid token"}
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		inner.ServeHTTP(w, r)

	})
}
