package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binder-project/binder-registry/registry"
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

func main() {
	apiKey := os.Getenv("BINDER_API_KEY")

	if apiKey == "" {
		fmt.Println("Environment variable BINDER_API_KEY must be non-empty")
		os.Exit(2)
	}

	store := registry.NewInMemoryStore()
	ctxt := NewRegistryContext(store, apiKey)

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
