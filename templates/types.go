package templates

import (
	"io"
)

const (
	// TypeGo is the type for Go templates.
	TypeGo = "go"

	// TypeHandlebars is the type for handlebars (mustache) templates
	TypeHandlebars = "handlebars"
)

// IFactory represents the available functions of the templates factory
type IFactory interface {

	// New creates a new default template type, defined in the configuration. If the default template type is not found, returns the implementation of Go Templates.
	New(name ...string) ITemplate

	// Create creates a template implementation by the given template type.
	Create(typeStr string, name ...string) (ITemplate, error)

	// Register registers a constructor for a template implementation.
	Register(name string, constructor func(name ...string) ITemplate)

	// GetAvaliableTypes returns the available type of template implementations
	GetAvaliableTypes() []string

	// Contains indicates whether an implementation of the given type is registered in the factory.
	Contains(typeStr string) bool
}

// ITemplate represents the templates interface to be used for template parsing.
type ITemplate interface {

	// Name represents the name of the ITemplate instance.
	Name() string

	// Type returns the template type of this instance.
	Type() string

	// ParseJSON parses the template using the string representation of a JSON
	ParseJSON(writer io.Writer, jsonStr string) error

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
}

type baseTemplate struct {
	ITemplate
	name string
}

func (b *baseTemplate) Type() string {
	return b.name
}
