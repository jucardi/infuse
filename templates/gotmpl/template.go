package gotmpl

import (
	"fmt"
	"io"
	"text/template"

	"github.com/jucardi/go-infuse/templates/base"
	"github.com/jucardi/go-strings"
)

// GoTemplate represents the implementation of ITemplate for Go templates
type GoTemplate struct {
	*base.AbstractTemplate
}

// Type returns the template type of this instance.
func (t *GoTemplate) Type() string {
	return "go"
}

// Parse parses the template
func (t *GoTemplate) Parse(writer io.Writer, data interface{}) error {
	str := t.prepare()
	tmpl := template.New(t.NameStr).Funcs(defaultFuncMap())
	if _, err := tmpl.Parse(str); err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

// LoadTemplate loads the given string as the template to be parsed.
func (t *GoTemplate) LoadTemplate(tmpl string) error {
	return validate(t.NameStr, tmpl, func() {
		t.Template = tmpl
	})
}

// LoadDefinition loads the give template string as a definition {{define "name"}}, using the given name as the name of the definition, to be used for 'template' directives.
func (t *GoTemplate) LoadDefinition(name, tmpl string) error {
	return validate(name, tmpl, func() {
		t.Definitions[name] = tmpl
	})
}

func (t *GoTemplate) prepare() string {
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
func New(name ...string) *GoTemplate {
	gt := &GoTemplate{}
	bt := &base.AbstractTemplate{
		IAbstractTemplateMembers: gt,
		NameStr:                  stringx.GetOrDefault("base", name...),
		Definitions:              map[string]string{},
	}
	gt.AbstractTemplate = bt
	return gt
}

func validate(name, tmpl string, successFn func()) error {
	_, err := template.New(name).Funcs(defaultFuncMap()).Parse(tmpl)

	if err != nil {
		return fmt.Errorf("unable to load definition '%s', %v", name, err)
	}

	successFn()
	return nil
}
