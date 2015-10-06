package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binder-project/binder-registry/registry"
)

func main() {
	apiKey := os.Getenv("BINDER_API_KEY")

	if apiKey == "" {
		fmt.Println("Environment variable BINDER_API_KEY must be non-empty")
		os.Exit(2)
	}

	store := registry.NewInMemoryStore()
	authStore := registry.NewTokenAuthStore(apiKey)

	reg := registry.Registry{store, authStore}
	router := registry.NewDefaultRouter(reg)

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
