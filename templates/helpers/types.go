package helpers

// IHelpersManager represents a template helpers managers
type IHelpersManager interface {
	// Get returns the helpers that have been registered to the manager
	Get() []*Helper

	// Register registers a helper function to be used with Go templates.
	Register(name string, fn interface{}, description ...string) error

	// Contains indicates whether a helper is contained by the manager
	Contains(name string) bool

	// Remove removes a helper by the given name, does nothing if the helper is not present
	Remove(name string)
}

// Helper encapsulates the information of a template helper
type Helper struct {
	Name        string
	Category    string
	Description string
	Function    interface{}
}
