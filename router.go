package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter sets up a mux.Router with the registry routes
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		// Wrap our routes with logging
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
