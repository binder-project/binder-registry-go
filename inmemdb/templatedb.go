package inmemdb

// inmemdb is a non-thread-safe in-memory database for the moment

import (
	"errors"
	"time"

	"github.com/binder-project/binder-registry/registry"
)

// InMemoryStore implements a non-thread-safe registry of binder templates
type InMemoryStore struct {
	templateMap map[string]registry.Template
}

// NewInMemoryStore create a new InMemoryStore
func NewInMemoryStore() InMemoryStore {
	var store InMemoryStore
	store.templateMap = make(map[string]registry.Template)
	return store
}

// GetTemplate retrieves the template with name, erroring otherwise
func (store InMemoryStore) GetTemplate(name string) (registry.Template, error) {
	tmpl, ok := store.templateMap[name]
	if !ok {
		return registry.Template{}, errors.New("Template unavailable")
	}

	return tmpl, nil
}

// RegisterTemplate registers the template in the Store
func (store InMemoryStore) RegisterTemplate(tmpl registry.Template) (registry.Template, error) {
	// Ensure tmpl.Name is available
	_, exists := store.templateMap[tmpl.Name]
	if exists {
		// Disallow registration if it exists
		return registry.Template{}, errors.New("Template already exists")
	}

	// Apply creation times
	tmpl.TimeModified = time.Now().UTC()
	tmpl.TimeCreated = tmpl.TimeModified

	store.templateMap[tmpl.Name] = tmpl
	return tmpl, nil
}

// ListTemplates provides a listing of all registered templates
func (store InMemoryStore) ListTemplates() ([]registry.Template, error) {
	templates := make([]registry.Template, len(store.templateMap))
	i := 0
	for _, tmpl := range store.templateMap {
		templates[i] = tmpl
		i++
	}
	return templates, nil
}

// UpdateTemplate will allow for updating ImageName and Command
func (store InMemoryStore) UpdateTemplate(tmpl registry.Template) (registry.Template, error) {
	updatedTemplate, ok := store.templateMap[tmpl.Name]
	if !ok {
		return registry.Template{}, errors.New("Template unavailable")
	}

	// For now we allow updates to image name and command
	if tmpl.ImageName != "" {
		updatedTemplate.ImageName = tmpl.ImageName
	}
	if tmpl.Command != "" {
		updatedTemplate.Command = tmpl.Command
	}

	// TODO: If fields are set inappropriately, return new error
	updatedTemplate.TimeModified = time.Now().UTC()

	store.templateMap[tmpl.Name] = updatedTemplate

	return updatedTemplate, nil
}
