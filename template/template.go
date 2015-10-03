package template

import "time"

// ContainerLimits is the specific set of limits we have for containers
// TODO: Should this come from Docker itself?
type ContainerLimits struct {
	Memory    int64 `json:"memory,omitempty"`
	CPUShares int64 `json:"cpu_shares,omitempty"`
}

// Template defines an image to be run along with its configuration
type Template struct {
	// Likely addition to spec
	Name string `json:"name"`

	ImageName string `json:"image-name"`
	Command   string `json:"command"`

	Limits ContainerLimits `json:"limits,omitempty"`

	// TODO: Add these to the binder API spec?
	RedirectURI string `json:"redirect-uri,omitempty"`
	BindIP      string `json:"container-ip,omitempty"`
	BindPort    int64  `json:"bind-port,omitempty"`

	// TODO: Also for the binder API spec, temporal understanding
	//       AKA good ol' timestamps
	TimeCreated  time.Time `json:"time-created"`
	TimeModified time.Time `json:"time-modified"`
}

// RegistryDB is an interface for registering templates into some backend
// storage, which could be in-memory, mongo, Postgres, Bolt, etc.
//
// This current interface makes no guarantees about thread safety, that's
// up to the implementer!
type RegistryDB interface {
	// GetTemplate retrieves the template with name, erroring otherwise
	GetTemplate(name string) (Template, error)

	// RegisterTemplate registers the template in the DB
	RegisterTemplate(tmpl Template) (Template, error)

	// ListTemplates provides a listing of all registered templates
	ListTemplates() ([]Template, error)

	// UpdateTemplate will allow for updating ImageName and Command
	UpdateTemplate(tmpl Template) (Template, error)
}
