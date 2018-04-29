package handlebars

import (
	"github.com/aymerick/raymond"
	"github.com/jucardi/go-infuse/templates/helpers"
)

type hbHelpersManager struct {
	helpers.IHelpersManager
}

// Register registers a helper function to be used with Go templates.
func (h *hbHelpersManager) Register(name string, fn interface{}, description ...string) error {
	if err := h.IHelpersManager.Register(name, fn, description...); err != nil {
		return err
	}
	raymond.RegisterHelper(name, fn)
	return nil
}

// Remove removes a helper by the given name, does nothing if the helper is not present
func (h *hbHelpersManager) Remove(name string) {

}

var instance *hbHelpersManager

// Helpers returns the singleton helpers instance used for Go templates
func Helpers() helpers.IHelpersManager {
	if instance == nil {
		base := helpers.New()
		instance = &hbHelpersManager{
			IHelpersManager: base,
		}
	}
	return instance
}

func init() {
	helpers.RegisterCommon(Helpers())
}
