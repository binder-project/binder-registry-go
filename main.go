package main

import (
	"log"
	"net/http"

	"github.com/binder-project/binder-registry/registry"
)

// RegistryContext keeps context for Store with the API Handlers
type RegistryContext struct {
	registry.Store
}

// NewRegistryContext initializes the context with a backend
func NewRegistryContext(store registry.Store) RegistryContext {
	return RegistryContext{
		Store: store,
	}
}

func main() {
	store := registry.NewInMemoryStore()
	ctxt := NewRegistryContext(store)

	var routes = []Route{
		Route{
			"Index",
			"GET",
			"/",
			ctxt.Index,
		},
		Route{
			"TemplateIndex",
			"GET",
			"/templates",
			ctxt.TemplateIndex,
		},
		Route{
			"TemplateShow",
			"GET",
			"/templates/{templateName}",
			ctxt.TemplateShow,
		},
		Route{
			"TemplateCreate",
			"POST",
			"/templates",
			ctxt.TemplateCreate,
		},
		Route{
			"TemplateUpdate",
			"PUT",
			"/templates/{templateName}",
			ctxt.TemplateUpdate,
		},
	}
	router := NewRouter(routes)

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Route is a simple HTTP Method, Pattern, and Handler pairing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
