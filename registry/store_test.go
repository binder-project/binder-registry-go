package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SUCCESSFUL TEST"))
}

func setup() {
	authStore = TokenAuthStore{"HOKEYPOKEY"}
	authorize = authStore.Authorize(http.HandlerFunc(EmptyHandler))
	req, _ = http.NewRequest("POST", "", nil)
	w = httptest.NewRecorder()
}

func teardown() {
	req = nil
	w = nil
}

func TestInvalidToken(t *testing.T) {
	setup()
	defer teardown()

	req.Header = map[string][]string{
		"Authorization": {"token BADKEY"},
	}
	authorize.ServeHTTP(w, req)

	matchedAPIError(t, InvalidTokenError, w)
}

func TestValidToken(t *testing.T) {
	setup()
	defer teardown()

	req.Header = map[string][]string{
		"Authorization": {"token HOKEYPOKEY"},
	}
	authorize.ServeHTTP(w, req)
	equals(t, "SUCCESSFUL TEST", w.Body.String())
}

func TestNoAuthHeader(t *testing.T) {
	setup()
	defer teardown()

	req.Header = map[string][]string{}
	authorize.ServeHTTP(w, req)

	matchedAPIError(t, MissingAuthHeaderError, w)
}

func MissingTokenField(t *testing.T) {
	setup()
	defer teardown()

	req.Header = map[string][]string{
		"Authorization": {"api-key HOKEYPOKEY"},
	}
	authorize.ServeHTTP(w, req)

	matchedAPIError(t, MissingTokenFieldError, w)
}
