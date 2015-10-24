package registry

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/justinas/alice"
)

// PanicStore only ever panics, used to make sure that responses to users
// still go out as application/json with appropriate errors, regardless of
// whether the backend store is causing fits. This can't perfectly handle the
// near heat death of the universe, but you'll at least have your towel.
type PanicStore struct {
}

func (s PanicStore) GetTemplate(name string) (Template, error) {
	panic("NO TEMPLATE FOR YOU")
}

func (s PanicStore) RegisterTemplate(tmpl Template) (Template, error) {
	panic("NO REGISTERING TEMPLATES FOR YOU")
}

func (s PanicStore) ListTemplates() ([]Template, error) {
	panic("NO LISTING TEMPLATES FOR YOU")
}

func (s PanicStore) UpdateTemplate(name string, update map[string]interface{}) (Template, error) {
	panic("NO UPDATING TEMPLATES FOR YOU")
}

func (s PanicStore) Authorize(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("NO AUTHORIZATION FOR YOU")
	})
}

type hokeyRequest struct {
	Request *http.Request
}

func newHokeyRequest(method string, urlStr string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, urlStr, body)
	return req
}

func TestPanicRecovery(t *testing.T) {
	registry := Registry{
		Store:     PanicStore{},
		AuthStore: PanicStore{},
	}

	handlerFuncs := []http.HandlerFunc{
		registry.TemplateCreate,
		registry.TemplateUpdate,
		registry.TemplateShow,
		registry.TemplateIndex,
	}

	middleware := alice.New(recoverHandler)

	for _, handlerFunc := range handlerFuncs {
		req, _ := http.NewRequest("GET", "", nil)
		w := httptest.NewRecorder()
		middleware.Then(handlerFunc).ServeHTTP(w, req)

		matchedAPIError(t, DontPanicError, w)
	}

	// Panics should be recovered from with the DefaultRouter
	router := NewDefaultRouter(registry)
	requests := []*http.Request{
		newHokeyRequest("GET", "/templates", strings.NewReader("{}")),
		newHokeyRequest("POST", "/templates", strings.NewReader("{}")),
		newHokeyRequest("GET", "/templates/env", strings.NewReader("{}")),
		newHokeyRequest("PUT", "/templates/env", strings.NewReader("{}")),
	}

	for _, req := range requests {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		matchedAPIError(t, DontPanicError, w)
	}

}
