package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"

	"github.com/binder-project/binder-registry/registry"
	"github.com/gorilla/mux"
)

// RegistryContext keeps context for Store with the API Handlers
type RegistryContext struct {
	registry.Store
	Token string
	Name  string
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

	common := alice.New(logger, contentTypeJSON)

	index := common.ThenFunc(ctxt.Index)
	templateIndex := common.ThenFunc(ctxt.TemplateIndex)
	templateShow := common.ThenFunc(ctxt.TemplateShow)

	authed := common.Append(ctxt.Authenticated)

	templateCreate := authed.ThenFunc(ctxt.TemplateCreate)
	templateUpdate := authed.ThenFunc(ctxt.TemplateUpdate)

	router.Methods("GET").Path("/").Handler(index)
	router.Methods("GET").Path("/templates").Handler(templateIndex)
	router.Methods("POST").Path("/templates").Handler(templateCreate) // templateCreate)
	router.Methods("GET").Path("/templates/{templateName}").Handler(templateShow)
	router.Methods("PUT").Path("/templates/{templateName}").Handler(templateUpdate)

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
