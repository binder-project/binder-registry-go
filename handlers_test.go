package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/binder-project/binder-registry/mocks"
	"github.com/gorilla/mux"
)

func setupHandlersTests() (Registry, *mocks.Store, *mocks.AuthStore) {
	store := new(mocks.Store)
	authStore := new(mocks.AuthStore)
	registry := Registry{
		Store:     store,
		AuthStore: authStore,
	}

	return registry, store, authStore
}

func TestIndex(t *testing.T) {
	registry, _, _ := setupHandlersTests()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	registry.Index(w, req)

	equals(t, "{\"status\": \"Binder Registry Live!\"}\n", w.Body.String())
}

func TestTemplateIndex(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mStore.On("ListTemplates").Return([]Template{}, nil)
	registry.TemplateIndex(w, req)
	mStore.AssertExpectations(t)
	contentTypeWasJSON(t, w)
	equals(t, w.Body.String(), "[]\n")
}

func TestTemplateIndexErrorHandle(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mStore.On("ListTemplates").Return([]Template{}, errors.New("Catastrophic"))
	registry.TemplateIndex(w, req)
	mStore.AssertExpectations(t)
	matchedAPIError(t, UnableToListError, w)
}

func TestTemplateCreate(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	myTemplate := Template{
		Name:      "test",
		ImageName: "jupyter/test",
	}

	body, err := json.Marshal(myTemplate)
	ok(t, err)

	b := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/templates/", b)
	w := httptest.NewRecorder()

	mStore.registerCall = make(map[string]RegisterReceiver)
	mStore.registerCall[myTemplate.Name] = RegisterReceiver{
		Input:          myTemplate,
		OutputTemplate: myTemplate,
		OutputError:    nil,
	}

	registry.TemplateCreate(w, req)

	reciever := mStore.registerCall[myTemplate.Name]
	equals(t, myTemplate.Name, reciever.Input.Name)
	equals(t, myTemplate.ImageName, reciever.Input.ImageName)
	equals(t, myTemplate.Name, reciever.Input.Name)
	equals(t, myTemplate.ImageName, reciever.Input.ImageName)
	ok(t, reciever.OutputError)
}

func TestTemplateShowNotFound(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	mStore.On("GetTemplate", "mytest").Return(Template{}, errors.New("Catastrophic"))

	req, _ := http.NewRequest("GET", "/templates/mytest", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/templates/{templateName}").HandlerFunc(registry.TemplateShow)
	router.ServeHTTP(w, req)
	mStore.AssertExpectations(t)
	matchedAPIError(t, TemplateNotFoundError, w)
}

func TestTemplateShow(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	myTemplate := Template{
		Name:      "mytest",
		ImageName: "jupyter/fun",
	}
	mStore.On("GetTemplate", "mytest").Return(myTemplate, nil)
	req, _ := http.NewRequest("GET", "/templates/mytest", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/templates/{templateName}").HandlerFunc(registry.TemplateShow)
	router.ServeHTTP(w, req)
	mStore.AssertExpectations(t)

	var actualTemplate Template
	json.NewDecoder(w.Body).Decode(&actualTemplate)

	equals(t, myTemplate.Name, actualTemplate.Name)
	equals(t, myTemplate.ImageName, actualTemplate.ImageName)
}

func TestTemplateUpdate(t *testing.T) {
	// This one, out of all of them, is testing nothing.
	registry, mStore, _ := setupHandlersTests()
	myTemplate := Template{
		Name:      "test",
		ImageName: "jupyter/test",
	}

	body, err := json.Marshal(myTemplate)
	ok(t, err)

	b := bytes.NewReader(body)
	req, _ := http.NewRequest("PUT", "/templates/test", b)
	w := httptest.NewRecorder()

	mStore.registerCall = make(map[string]RegisterReceiver)
	mStore.registerCall[myTemplate.Name] = RegisterReceiver{
		Input:          myTemplate,
		OutputTemplate: myTemplate,
		OutputError:    nil,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("PUT").Path("/templates/{templateName}").HandlerFunc(registry.TemplateUpdate)
	router.ServeHTTP(w, req)

	reciever := mStore.registerCall[myTemplate.Name]
	equals(t, myTemplate.Name, reciever.Input.Name)
	equals(t, myTemplate.ImageName, reciever.Input.ImageName)
	equals(t, myTemplate.Name, reciever.Input.Name)
	equals(t, myTemplate.ImageName, reciever.Input.ImageName)
	ok(t, reciever.OutputError)
}
