package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/binder-project/binder-registry/template"

	"github.com/gorilla/mux"
)

// Index lists available resources at this endpoint
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Binder Registry Live!")
}

// TemplateIndex lists the available templates as well as their configuration
func TemplateIndex(w http.ResponseWriter, r *http.Request) {
	templates := ListTemplates()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(templates); err != nil {
		panic(err)
	}
}

// TemplateShow displays the requested template
func TemplateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl, err := GetTemplate(templateName)

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
func TemplateCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tmpl template.Template
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

	t, err := RegisterTemplate(tmpl)
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
func TemplateUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tmpl template.Template
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		// TODO: Return some minimal error back
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl.Name = templateName

	tmpl, err = UpdateTemplate(tmpl)
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
