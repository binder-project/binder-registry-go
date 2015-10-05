package registry

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
