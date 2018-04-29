package helpers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jucardi/go-strings"
)

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
	Description string
	Function    interface{}
}

// HelpersManager is the basic implementation of IHelpersManager, administers helpers to be used by templates
type HelpersManager struct {
	helpers map[string]*Helper
}

// Get returns the helpers that have been registered to the manager
func (a *HelpersManager) Get() []*Helper {
	var ret []*Helper
	for _, v := range a.helpers {
		ret = append(ret, v)
	}
	return ret
}

// Register registers a helper function to be used with Go templates.
func (a *HelpersManager) Register(name string, fn interface{}, description ...string) error {
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
func (a *HelpersManager) Contains(name string) bool {
	_, ok := a.helpers[name]
	return ok
}

// Remove removes a helper by the given name, does nothing if the helper is not present
func (a *HelpersManager) Remove(name string) {
	delete(a.helpers, name)
}

// New returns a new instance of HelpersManager
func New() IHelpersManager {
	return &HelpersManager{helpers: map[string]*Helper{}}
}
