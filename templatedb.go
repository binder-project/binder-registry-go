// templatedb is a non-thread-safe in-memory database for the moment
package main

import (
	"errors"

	"github.com/binder-project/binder-registry/template"
)

var templateMap map[string]template.Template

func init() {
	templateMap = make(map[string]template.Template)
}

// GetTemplate retrieves the template with name, erroring otherwise
func GetTemplate(name string) (template.Template, error) {
	tmpl, ok := templateMap[name]
	if !ok {
		return template.Template{}, errors.New("Template unavailable")
	}

	return tmpl, nil
}

// RegisterTemplate registers the template in the DB
func RegisterTemplate(tmpl template.Template) (template.Template, error) {
	// Ensure tmpl.Name is available
	_, exists := templateMap[tmpl.Name]
	if exists {
		// Disallow registration if it exists
		return template.Template{}, errors.New("Template already exists")
	}

	templateMap[tmpl.Name] = tmpl
	return tmpl, nil
}

// ListTemplates provides a listing of all registered templates
func ListTemplates() []template.Template {
	templates := make([]template.Template, len(templateMap))
	i := 0
	for _, tmpl := range templateMap {
		templates[i] = tmpl
		i++
	}
	return templates
}
