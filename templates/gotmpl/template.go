package gotmpl

import (
	"fmt"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/templates/base"
	"gopkg.in/jucardi/go-strings.v1/stringx"
	"io"
	"text/template"
)

// TypeGo is the type for Go templates.
const TypeGo = "go"

func init() {
	templates.Factory().Register(TypeGo, func(name ...string) templates.ITemplate { return New(name...) })
}

// Template represents the implementation of ITemplate for Go templates
type Template struct {
	*base.AbstractTemplate
}

// Type returns the template type of this instance.
func (t *Template) Type() string {
	return "go"
}

// Parse parses the template
func (t *Template) Parse(writer io.Writer, data interface{}) error {
	str := t.prepare()
	tmpl := template.New(t.NameStr).Funcs(getHelpers())
	if _, err := tmpl.Parse(str); err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

// LoadTemplate loads the given string as the template to be parsed.
func (t *Template) LoadTemplate(tmpl string) error {
	return validate(t.NameStr, tmpl, func() {
		t.Template = tmpl
	})
}

// LoadDefinition loads the give template string as a definition {{define "name"}}, using the given name as the name of the definition, to be used for 'template' directives.
func (t *Template) LoadDefinition(name, tmpl string) error {
	return validate(name, tmpl, func() {
		t.Definitions[name] = tmpl
	})
}

func (t *Template) prepare() string {
	builder := stringx.Builder()

	for k, v := range t.Definitions {
		if k == t.NameStr {
			continue
		}
		builder.
			AppendLinef("{{define \"%s\"}}", k).
			AppendLine(v).
			AppendLine("{{end}}")
	}
	return builder.AppendLine(t.Template).Build()
}

// New creates a new template utility which extends the default built in functions for Go templates.
func New(name ...string) *Template {
	gt := &Template{}
	bt := &base.AbstractTemplate{
		IAbstractTemplateMembers: gt,
		NameStr:                  stringx.GetOrDefault("base", name...),
		Definitions:              map[string]string{},
	}
	gt.AbstractTemplate = bt
	return gt
}

func validate(name, tmpl string, successFn func()) error {
	_, err := template.New(name).Funcs(getHelpers()).Parse(tmpl)

	if err != nil {
		return fmt.Errorf("unable to load definition '%s', %v", name, err)
	}

	successFn()
	return nil
}
