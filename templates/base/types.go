package base

import "io"

// IAbstractTemplateMembers represents the templates interface to be used for template parsing.
type IAbstractTemplateMembers interface {

	// Parse parses the template
	Parse(writer io.Writer, data interface{}) error

	// LoadTemplate loads the given string as the template to be parsed.
	LoadTemplate(tmpl string) error

	// LoadDefinition loads the given template string as a definition to be used for 'template' directives.
	LoadDefinition(name, tmpl string) error
}
