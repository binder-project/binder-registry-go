// templatedb is a non-thread-safe in-memory database for the moment
package main

import (
	"errors"
	"time"

	"github.com/binder-project/binder-registry/template"
)

var registry template.RegistryDB

// For the time being we're going to auto-init here
func init() {
	registry = NewInMemoryRegistryDB()
}

// InMemoryRegistryDB implements a non-thread-safe registry of binder templates
type InMemoryRegistryDB struct {
	templateMap map[string]template.Template
}

// NewInMemoryRegistryDB create a new InMemoryRegistryDB
func NewInMemoryRegistryDB() InMemoryRegistryDB {
	var registry InMemoryRegistryDB
	registry.templateMap = make(map[string]template.Template)
	return registry
}

// GetTemplate retrieves the template with name, erroring otherwise
func (reg InMemoryRegistryDB) GetTemplate(name string) (template.Template, error) {
	tmpl, ok := reg.templateMap[name]
	if !ok {
		return template.Template{}, errors.New("Template unavailable")
	}

	return tmpl, nil
}

// RegisterTemplate registers the template in the DB
func (reg InMemoryRegistryDB) RegisterTemplate(tmpl template.Template) (template.Template, error) {
	// Ensure tmpl.Name is available
	_, exists := reg.templateMap[tmpl.Name]
	if exists {
		// Disallow registration if it exists
		return template.Template{}, errors.New("Template already exists")
	}

	// Apply creation times
	tmpl.TimeModified = time.Now().UTC()
	tmpl.TimeCreated = tmpl.TimeModified

	reg.templateMap[tmpl.Name] = tmpl
	return tmpl, nil
}

// ListTemplates provides a listing of all registered templates
func (reg InMemoryRegistryDB) ListTemplates() ([]template.Template, error) {
	templates := make([]template.Template, len(reg.templateMap))
	i := 0
	for _, tmpl := range reg.templateMap {
		templates[i] = tmpl
		i++
	}
	return templates, nil
}

// UpdateTemplate will allow for updating ImageName and Command
func (reg InMemoryRegistryDB) UpdateTemplate(tmpl template.Template) (template.Template, error) {
	updatedTemplate, ok := reg.templateMap[tmpl.Name]
	if !ok {
		return template.Template{}, errors.New("Template unavailable")
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

	reg.templateMap[tmpl.Name] = updatedTemplate

	return updatedTemplate, nil
}
