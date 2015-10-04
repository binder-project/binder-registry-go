package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/binder-project/binder-registry/registry"

	"github.com/gorilla/mux"
)

// Index lists available resources at this endpoint
func (ctxt RegistryContext) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Binder Registry Live!")
}

// TemplateIndex lists the available templates as well as their configuration
func (ctxt RegistryContext) TemplateIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	templates, err := ctxt.ListTemplates()
	if err != nil {
		// TODO: Pick appropriate response code
		w.WriteHeader(500)
		userErr := jsonErr{Code: 500, Text: "Unable to list templates."}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(templates); err != nil {
		panic(err)
	}
}

// TemplateShow displays the requested template
func (ctxt RegistryContext) TemplateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl, err := ctxt.GetTemplate(templateName)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		userErr := jsonErr{Code: http.StatusNotFound, Text: "Template Not Found"}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
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
func (ctxt RegistryContext) TemplateCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tmpl registry.Template
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		userErr := jsonErr{Code: http.StatusBadRequest, Text: err.Error()}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	if tmpl.Name == "" || tmpl.ImageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		userErr := jsonErr{Code: http.StatusBadRequest, Text: "name and image-name must be specified"}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	t, err := ctxt.RegisterTemplate(tmpl)
	if err != nil {
		w.WriteHeader(400) // That or 409 Conflict
		userErr := jsonErr{Code: 400, Text: err.Error()}
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
func (ctxt RegistryContext) TemplateUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tmpl registry.Template
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		userErr := jsonErr{Code: 422, Text: err.Error()}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl.Name = templateName

	tmpl, err = ctxt.UpdateTemplate(tmpl)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		userErr := jsonErr{Code: http.StatusNotFound, Text: "Template Not Found"}
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
