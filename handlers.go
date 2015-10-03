package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

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
	vars := mux.Vars(r)
	templateName := vars["templateName"]

	tmpl, err := GetTemplate(templateName)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		userErr := jsonErr{Code: http.StatusNotFound, Text: "Template Not Found"}
		if err := json.NewEncoder(w).Encode(userErr); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	tmpl.TimeCreated = time.Now().UTC()
	tmpl.TimeModified = tmpl.TimeCreated
	t := CreateTemplate(tmpl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
