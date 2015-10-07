package registry

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Store is an interface for registering templates into some backend
// storage, which could be in-memory, mongo, Postgres, Bolt, etc.
//
// This current interface makes no guarantees about thread safety, that's
// up to the implementer of the interface!
type Store interface {
	// GetTemplate retrieves the template with name, erroring otherwise
	GetTemplate(name string) (Template, error)

	// RegisterTemplate registers the template in the DB
	RegisterTemplate(tmpl Template) (Template, error)

	// ListTemplates provides a listing of all registered templates
	ListTemplates() ([]Template, error)

	// UpdateTemplate will allow for updating ImageName and Command
	UpdateTemplate(tmpl Template) (Template, error)
}

// AuthStore is an interface for connecting to some authentication endpoint,
// assumed to be used as middleware in front of authenticated endpoints
// much like the default setup for TemplateUpdate or TemplateCreate
type AuthStore interface {
	Authorize(inner http.Handler) http.Handler
}

// TokenAuthStore uses a single token for authentication
type TokenAuthStore struct {
	Token string
}

// Authorize uses token based authentication, presuming that the endpoint is
// done over HTTPS
func (authStore TokenAuthStore) Authorize(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		authHeader, hasAuthHeader := r.Header["Authorization"]

		if !hasAuthHeader || len(authHeader) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(MissingAuthHeaderError); err != nil {
				panic(err)
			}
			return
		}

		// QUESTION: What happens when there are multiple Authorization headers?
		tokenLine := authHeader[0]
		authHeader = strings.SplitN(tokenLine, " ", 3)

		if len(authHeader) != 2 || strings.ToLower(authHeader[0]) != "token" {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(MissingTokenFieldError); err != nil {
				panic(err)
			}
			return
		}

		putativeToken := authHeader[1]

		// TODO: HMAC SHA between shared secret + values, minding timing attacks
		// TODO #2: See if go has some of this built in nicely

		if putativeToken != authStore.Token {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(InvalidTokenError); err != nil {
				panic(err)
			}
			return
		}

		inner.ServeHTTP(w, r)

	})
}

// InvalidTokenError is an APIErrorResponse for when a token is invalid
var InvalidTokenError = APIErrorResponse{Message: "Invalid token"}

// MissingTokenFieldError is an APIErrorResponse for when the token field is missing
var MissingTokenFieldError = APIErrorResponse{
	Message: "Token field not present in Authorization header. Should be of format 'Authorization: token key'",
}

// MissingAuthHeaderError is an APIErrorResponse for when the Authorization header is not set
var MissingAuthHeaderError = APIErrorResponse{
	Message: "Authorization header not set. Should be of format 'Authorization: token key'",
}
