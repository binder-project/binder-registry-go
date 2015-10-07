package registry

import (
	"time"
)

// ContainerLimits is the specific set of limits we have for containers
// QUESTION: Should this come from Docker itself?
type ContainerLimits struct {
	Memory    int64 `json:"memory,omitempty"`
	CPUShares int64 `json:"cpu_shares,omitempty"`
}

// Template defines an image to be run along with its configuration
type Template struct {
	Name string `json:"name"`

	ImageName string `json:"image-name"`
	Command   string `json:"command"`

	Limits ContainerLimits `json:"limits,omitempty"`

	TimeCreated  time.Time `json:"time-created"`
	TimeModified time.Time `json:"time-modified"`

	// TODO: Add these to the binder API spec?
	RedirectURI string `json:"redirect-uri,omitempty"`
	BindIP      string `json:"container-ip,omitempty"`
	BindPort    int64  `json:"bind-port,omitempty"`
}
