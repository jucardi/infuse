package handlebars

import (
	"github.com/aymerick/raymond"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/templates/base"
	"github.com/jucardi/infuse/templates/helpers"
	"io"
)

// TypeHandlebars is the type for handlebars (mustache) templates
const TypeHandlebars = "handlebars"

func init() {
	templates.Factory().Register(TypeHandlebars, func(name ...string) templates.ITemplate { return New(name...) })
}

// Template represents the implementation of ITemplate for handlebars (mustache) templates
type Template struct {
	*base.AbstractTemplate
}

// Type returns the template type of this instance.
func (t *Template) Type() string {
	return "handlebars"
}

// Parse parses the template
func (t *Template) Parse(writer io.Writer, data interface{}) error {
	str, err := raymond.Render(t.Template, data)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(str))
	return err
}

// LoadTemplate loads the given string as the template to be parsed.
func (t *Template) LoadTemplate(tmpl string) error {
	t.Template = tmpl
	return nil
}

// LoadDefinition loads the given template string as a definition to be used for 'template' directives.
func (t *Template) LoadDefinition(name, tmpl string) error {
	t.Definitions[name] = tmpl
	return nil
}

func (t *Template) Helpers() (ret []*helpers.Helper) {
	return Helpers().Get()
}

// New creates a new template utility which extends the default built in functions for Go templates.
func New(name ...string) *Template {
	hb := &Template{}
	bt := &base.AbstractTemplate{
		IAbstractTemplateMembers: hb,
		NameStr:                  stringx.GetOrDefault("base", name...),
		Definitions:              map[string]string{},
	}
	hb.AbstractTemplate = bt
	return hb
}
