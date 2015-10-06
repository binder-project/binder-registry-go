package registry

// Registry keeps context for Store with the API Handlers
type Registry struct {
	Store
	AuthStore
}
