package registry

// inmemstore is a non-thread-safe in-memory database for the moment

import "time"

// InMemoryStore implements a non-thread-safe registry of binder templates
type InMemoryStore struct {
	templateMap map[string]Template
}

// NewInMemoryStore create a new InMemoryStore
func NewInMemoryStore() InMemoryStore {
	var store InMemoryStore
	store.templateMap = make(map[string]Template)
	return store
}

// GetTemplate retrieves the template with name, erroring otherwise
func (store InMemoryStore) GetTemplate(name string) (Template, error) {
	tmpl, ok := store.templateMap[name]
	if !ok {
		return Template{}, unavailableTemplateError
	}

	return tmpl, nil
}

// RegisterTemplate registers the template in the Store
func (store InMemoryStore) RegisterTemplate(tmpl Template) (Template, error) {
	// Ensure tmpl.Name is available
	_, exists := store.templateMap[tmpl.Name]
	if exists {
		// Disallow registration if it exists
		return Template{}, existingTemplateError
	}

	// Apply creation times
	tmpl.TimeModified = time.Now().UTC()
	tmpl.TimeCreated = tmpl.TimeModified

	store.templateMap[tmpl.Name] = tmpl
	return tmpl, nil
}

// ListTemplates provides a listing of all registered templates
func (store InMemoryStore) ListTemplates() ([]Template, error) {
	templates := make([]Template, len(store.templateMap))
	i := 0
	for _, tmpl := range store.templateMap {
		templates[i] = tmpl
		i++
	}
	return templates, nil
}

// UpdateTemplate will allow for updating ImageName and Command
func (store InMemoryStore) UpdateTemplate(tmpl Template) (Template, error) {
	updatedTemplate, ok := store.templateMap[tmpl.Name]
	if !ok {
		return Template{}, unavailableTemplateError
	}

	// For now we allow updates to image name and command
	if tmpl.ImageName != "" {
		updatedTemplate.ImageName = tmpl.ImageName
	}
	if tmpl.Command != "" {
		updatedTemplate.Command = tmpl.Command
	}

	updatedTemplate.TimeModified = time.Now().UTC()
	store.templateMap[tmpl.Name] = updatedTemplate

	return updatedTemplate, nil
}
