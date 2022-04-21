package helpers

import (
	"errors"
	"fmt"
	"github.com/jucardi/go-strings/stringx"
	"reflect"
)

// Manager is the basic implementation of IHelpersManager, administers helpers to be used by templates
type Manager struct {
	helpers map[string]*Helper
}

// Get returns the helpers that have been registered to the manager
func (a *Manager) Get() []*Helper {
	var ret []*Helper
	for _, v := range a.helpers {
		ret = append(ret, v)
	}
	return ret
}

// Register registers a helper function to be used with Go templates.
func (a *Manager) Register(name string, fn interface{}, description ...string) error {
	if name == "" {
		return errors.New("the helper name cannot be empty")
	}

	if fn == nil {
		return errors.New("the helper function cannot be nil")
	}

	if kind := reflect.TypeOf(fn).Kind(); kind != reflect.Func {
		return fmt.Errorf("wrong type for 'fn', %v , must be a function", kind)
	}

	a.helpers[name] = &Helper{
		Name:        name,
		Function:    fn,
		Description: stringx.GetOrDefault("", description...),
	}

	return nil
}

// Contains indicates whether a helper is contained by the manager
func (a *Manager) Contains(name string) bool {
	_, ok := a.helpers[name]
	return ok
}

// Remove removes a helper by the given name, does nothing if the helper is not present
func (a *Manager) Remove(name string) {
	delete(a.helpers, name)
}

// New returns a new instance of HelpersManager
func New() IHelpersManager {
	return &Manager{helpers: map[string]*Helper{}}
}
