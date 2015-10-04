package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binder-project/binder-registry/registry"
	"github.com/gorilla/mux"
)

// RegistryContext keeps context for Store with the API Handlers
type RegistryContext struct {
	registry.Store
	Token string
}

// NewRegistryContext initializes the context with a backend
func NewRegistryContext(store registry.Store, token string) RegistryContext {
	return RegistryContext{
		Store: store,
		Token: token,
	}
}

// NewRouter sets up a mux.Router with the registry routes
func NewRouter(ctxt RegistryContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").HandlerFunc(ctxt.Index)
	router.Methods("GET").Path("/templates").HandlerFunc(ctxt.TemplateIndex)
	router.Methods("GET").Path("/templates/{templateName}").HandlerFunc(ctxt.TemplateShow)

	router.Methods("POST").Path("/templates").HandlerFunc(ctxt.TemplateCreate)
	router.Methods("PUT").Path("/templates/{templateName}").HandlerFunc(ctxt.TemplateUpdate)
	return router
}

func main() {
	apiKey := os.Getenv("BINDER_API_KEY")

	if apiKey == "" {
		fmt.Println("Environment variable BINDER_API_KEY must be non-empty")
		os.Exit(2)
	}

	store := registry.NewInMemoryStore()
	ctxt := NewRegistryContext(store, apiKey)

	router := NewRouter(ctxt)

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
