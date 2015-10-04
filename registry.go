package main

import (
	"log"
	"net/http"

	"github.com/binder-project/binder-registry/inmemdb"
	"github.com/binder-project/binder-registry/registry"
)

var store registry.Store

func init() {
	// TODO: Switch to application context, a la
	//       https://gist.github.com/elithrar/5aef354a54ba71a32e23
	store = inmemdb.NewInMemoryStore()
}

func main() {

	router := NewRouter()

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
