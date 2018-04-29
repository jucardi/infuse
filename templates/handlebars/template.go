package handlebars

import (
	"io"

	"github.com/aymerick/raymond"
	"github.com/jucardi/go-infuse/templates/base"
	"github.com/jucardi/go-strings"
)

// HandlebarsTemplate represents the implementation of ITemplate for handlebars (mustache) templates
type HandlebarsTemplate struct {
	*base.AbstractTemplate
}

// Type returns the template type of this instance.
func (h *HandlebarsTemplate) Type() string {
	return "handlebars"
}

// Parse parses the template
func (h *HandlebarsTemplate) Parse(writer io.Writer, data interface{}) error {
	str, err := raymond.Render(h.Template, data)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(str))
	return err
}

// LoadTemplate loads the given string as the template to be parsed.
func (h *HandlebarsTemplate) LoadTemplate(tmpl string) error {
	h.Template = tmpl
	return nil
}

// LoadDefinition loads the given template string as a definition to be used for 'template' directives.
func (h *HandlebarsTemplate) LoadDefinition(name, tmpl string) error {
	h.Definitions[name] = tmpl
	return nil
}

// New creates a new template utility which extends the default built in functions for Go templates.
func New(name ...string) *HandlebarsTemplate {
	hb := &HandlebarsTemplate{}
	bt := &base.AbstractTemplate{
		IAbstractTemplateMembers: hb,
		NameStr:                  stringx.GetOrDefault("base", name...),
		Definitions:              map[string]string{},
	}
	hb.AbstractTemplate = bt
	return hb
}
