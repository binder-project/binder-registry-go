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
	Name string `bson:"name" json:"name"`

	ImageName string `bson:"image-name" json:"image-name"`
	Command   string `bson:"command" json:"command"`

	Limits ContainerLimits `bson:"limits" json:"limits,omitempty"`

	TimeCreated  time.Time `bson:"time-created" json:"time-created"`
	TimeModified time.Time `bson:"time-modified" json:"time-modified"`

	// TODO: Add these to the binder API spec?
	RedirectURI string `bson:"redirect-uri" json:"redirect-uri,omitempty"`
	BindIP      string `bson:"container-ip" json:"container-ip,omitempty"`
	BindPort    int64  `bson:"bind-port" json:"bind-port,omitempty"`
}
