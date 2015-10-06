package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// NewDefaultRouter sets up a mux.Router with the registry routes
func NewDefaultRouter(registry Registry) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	middleware := alice.New(recoverHandler, logger)

	index := middleware.ThenFunc(registry.Index)
	templateIndex := middleware.ThenFunc(registry.TemplateIndex)
	templateShow := middleware.ThenFunc(registry.TemplateShow)

	authed := middleware.Append(registry.AuthStore.Authorize)

	templateCreate := authed.ThenFunc(registry.TemplateCreate)
	templateUpdate := authed.ThenFunc(registry.TemplateUpdate)

	router.Methods("GET").Path("/").Handler(index)
	router.Methods("GET").Path("/templates").Handler(templateIndex)
	router.Methods("POST").Path("/templates").Handler(templateCreate) // templateCreate)
	router.Methods("GET").Path("/templates/{templateName}").Handler(templateShow)
	router.Methods("PUT").Path("/templates/{templateName}").Handler(templateUpdate)

	return router
}

var common alice.Chain

func init() {
	common = alice.New(contentTypeJSON)
}

// Index lists available resources at this endpoint
func (registry Registry) Index(w http.ResponseWriter, r *http.Request) {
	common.ThenFunc(registry.index).ServeHTTP(w, r)
}

func (registry Registry) index(w http.ResponseWriter, r *http.Request) {
	// TODO: Navigable API from the root
	fmt.Fprintln(w, "{\"status\": \"Binder Registry Live!\"}")
}

// TemplateIndex lists the available templates as well as their configuration
func (registry Registry) TemplateIndex(w http.ResponseWriter, r *http.Request) {
	common.ThenFunc(registry.templateIndex).ServeHTTP(w, r)
}

func (registry Registry) templateIndex(w http.ResponseWriter, r *http.Request) {
	templates, err := registry.ListTemplates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := APIErrorResponse{Message: "Unable to list templates"}
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(templates); err != nil {
		panic(err)
	}
}

// TemplateShow displays the requested template
func (registry Registry) TemplateShow(w http.ResponseWriter, r *http.Request) {
	common.ThenFunc(registry.templateShow).ServeHTTP(w, r)
}

func (registry Registry) templateShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl, err := registry.GetTemplate(templateName)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := APIErrorResponse{Message: "Template Not Found"}
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tmpl); err != nil {
		panic(err)
	}
}

// TemplateCreate registers a template
func (registry Registry) TemplateCreate(w http.ResponseWriter, r *http.Request) {
	common.ThenFunc(registry.templateCreate).ServeHTTP(w, r)
}

func (registry Registry) templateCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var tmpl Template

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		userErr := APIErrorResponse{Message: err.Error()}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	if tmpl.Name == "" || tmpl.ImageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		userErr := APIErrorResponse{Message: "name and image-name must be specified"}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	t, err := registry.RegisterTemplate(tmpl)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		userErr := APIErrorResponse{Message: err.Error()}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

// TemplateUpdate updates an individual template by name
func (registry Registry) TemplateUpdate(w http.ResponseWriter, r *http.Request) {
	common.ThenFunc(registry.templateUpdate).ServeHTTP(w, r)
}

func (registry Registry) templateUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var tmpl Template

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		userErr := APIErrorResponse{Message: err.Error()}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl.Name = templateName

	tmpl, err = registry.UpdateTemplate(tmpl)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		userErr := APIErrorResponse{Message: "Template Not Found"}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(tmpl); err != nil {
		panic(err)
	}
}
