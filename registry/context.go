package registry

// Registry keeps context for Store with the API Handlers
type Registry struct {
	Store
	AuthStore
	Name string
}

// NewRegistry initializes the context with a backend
func NewRegistry(store Store, authStore AuthStore) Registry {
	return Registry{
		Store:     store,
		AuthStore: authStore,
	}
}
