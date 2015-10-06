package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/binder-project/binder-registry/registry/mocks"
)

func setupHandlersTests() (Registry, *mockStore, *mocks.AuthStore) {
	store := new(mockStore)
	authStore := new(mocks.AuthStore)
	registry := Registry{
		Store:     store,
		AuthStore: authStore,
	}

	return registry, store, authStore
}

func TestIndex(t *testing.T) {
	registry, _, _ := setupHandlersTests()

	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	registry.Index(w, req)

	equals(t, "{\"status\": \"Binder Registry Live!\"}\n", w.Body.String())
}

func TestTemplateIndex(t *testing.T) {
	registry, mStore, _ := setupHandlersTests()
	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	mStore.On("ListTemplates").Return([]Template{}, nil)
	registry.TemplateIndex(w, req)
	mStore.AssertExpectations(t)

}
