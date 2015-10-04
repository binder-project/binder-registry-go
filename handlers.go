package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
		w.WriteHeader(http.StatusInternalServerError)
		userErr := jsonErr{Code: http.StatusInternalServerError, Text: "Unable to list templates."}
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

// Authorize probes for an Authorization header that matches the RegistryContext
func (ctxt RegistryContext) Authorize(w http.ResponseWriter, r *http.Request) (int, error) {
	authHeader, hasAuthHeader := r.Header["Authorization"]

	if !hasAuthHeader || len(authHeader) < 1 {
		return http.StatusUnauthorized, errors.New("Authorization header not set. Should be of format 'Authorization: token <key>'")
	}

	// QUESTION: What happens when there are multiple Authorization headers?
	tokenLine := authHeader[0]
	authHeader = strings.SplitN(tokenLine, " ", 3)
	if len(authHeader) != 2 || strings.ToLower(authHeader[0]) != "token" {
		return http.StatusUnauthorized, errors.New("Token field not present in Authorization header. Should be of format 'Authorization: token <key>'")
	}

	putativeToken := authHeader[1]

	// TODO: HMAC SHA between shared secret + values, minding timing attacks
	// TODO #2: See if go has some of this built in nicely

	if putativeToken != ctxt.Token {
		return http.StatusUnauthorized, errors.New("Invalid token")
	}

	return http.StatusAccepted, nil
}

// TemplateCreate registers a template
func (ctxt RegistryContext) TemplateCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	code, err := ctxt.Authorize(w, r)
	if err != nil {
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(jsonErr{code, err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	var tmpl registry.Template

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
		w.WriteHeader(http.StatusConflict)
		userErr := jsonErr{Code: http.StatusConflict, Text: err.Error()}
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

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	code, err := ctxt.Authorize(w, r)
	if err != nil {
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(jsonErr{code, err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	var tmpl registry.Template

	if err := json.Unmarshal(body, &tmpl); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		userErr := jsonErr{Code: http.StatusBadRequest, Text: err.Error()}
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
