package templates

import (
	"github.com/jucardi/infuse/templates/helpers"
	"io"
)

// IFactory represents the available functions of the templates factory
type IFactory interface {
	// New creates a new default template type, defined in the configuration. If the default template type is not found, returns the implementation of Go Templates.
	New(name ...string) ITemplate

	// Create creates a template implementation by the given template type.
	Create(typeStr string, name ...string) (ITemplate, error)

	// Register registers a constructor for a template implementation.
	Register(name string, constructor func(name ...string) ITemplate)

	// GetAvailableTypes returns the available type of template implementations
	GetAvailableTypes() []string

	// Contains indicates whether an implementation of the given type is registered in the factory.
	Contains(typeStr string) bool
}

// ITemplate represents the templates interface to be used for template parsing.
type ITemplate interface {
	// Name represents the name of the ITemplate instance.
	Name() string

	// Type returns the template type of this instance.
	Type() string

	// ParseMarshaled attempts to unmarshall the byte data provided and parses the template.
	ParseMarshaled(writer io.Writer, data []byte) error

	// ParseJSON attempts to unmarshall the byte data provided as a JSON object and parses the template.
	ParseJSON(writer io.Writer, data []byte) error

	// ParseYAML attempts to unmarshall the byte data provided as a YAML object and parses the template.
	ParseYAML(writer io.Writer, data []byte) error

	// Parse parses the template
	Parse(writer io.Writer, data interface{}) error

	// LoadFileTemplate loads the given file as the template to be parsed.
	LoadFileTemplate(filename string) error

	// LoadTemplate loads the given string as the template to be parsed.
	LoadTemplate(tmpl string) error

	// LoadFileDefinitionsByPattern uses a pattern to find the file definitions to be loaded for the template parsing
	LoadFileDefinitionsByPattern(pattern string) error

	// LoadFileDefinition loads a file(s) as definition(s) to be used for 'template' directives
	LoadFileDefinition(files ...string) error

	// LoadDefinition loads the given template string as a definition to be used for 'template' directives.
	LoadDefinition(name, tmpl string) error

	// Helpers returns the list of helpers that have been registered to this template
	Helpers() []*helpers.Helper
}

type baseTemplate struct {
	ITemplate
	name string
}

func (b *baseTemplate) Type() string {
	return b.name
}
