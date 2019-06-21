package gotmpl

import (
	"fmt"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/templates/base"
	"github.com/jucardi/infuse/templates/helpers"
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

func (t *Template) Helpers() (ret []*helpers.Helper) {
	ret = []*helpers.Helper{
		{Category: "Built-in Comparisons", Name: "eq", Description: "equals, evaluates the comparison a == b || a == c || ..."},
		{Category: "Built-in Comparisons", Name: "ge", Description: "greater than or equals, evaluates the comparison a >= b."},
		{Category: "Built-in Comparisons", Name: "gt", Description: "greater than, evaluates the comparison a > b."},
		{Category: "Built-in Comparisons", Name: "le", Description: "lower than or equals, evaluates the comparison <= b."},
		{Category: "Built-in Comparisons", Name: "lt", Description: "lower than, evaluates the comparison a < b."},
		{Category: "Built-in Comparisons", Name: "ne", Description: "not equals, evaluates the comparison a != b."},
		{Category: "Built-in Functions", Name: "and", Description: "computes the Boolean AND of its arguments, returning the first false argument it encounters, or the last argument."},
		{Category: "Built-in Functions", Name: "call", Description: "call returns the result of evaluating the first argument as a function. The function must return 1 result, or 2 results, the second of which is an error."},
		{Category: "Built-in Functions", Name: "html", Description: "returns the escaped HTML equivalent of the textual representation of its arguments."},
		{Category: "Built-in Functions", Name: "index", Description: "returns the result of indexing its first argument by the following arguments. Thus \"index x 1 2 3\" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array."},
		{Category: "Built-in Functions", Name: "js", Description: "returns the escaped JavaScript equivalent of the textual representation of its arguments"},
		{Category: "Built-in Functions", Name: "len", Description: "returns the length of the item, with an error if it has no defined length."},
		{Category: "Built-in Functions", Name: "not", Description: "returns the Boolean negation of its argument."},
		{Category: "Built-in Functions", Name: "or", Description: "computes the Boolean OR of its arguments, returning the first true argument it encounters, or the last argument."},
		{Category: "Built-in Functions", Name: "print", Description: "formats using the default formats for its operands and returns the resulting string. Spaces are added between operands when neither is a string."},
		{Category: "Built-in Functions", Name: "printf", Description: "formats according to a format specifier and returns the resulting string."},
		{Category: "Built-in Functions", Name: "println", Description: "formats using the default formats for its operands and returns the resulting string. Spaces are always added between operands and a newline is appended."},
		{Category: "Built-in Functions", Name: "urlquery", Description: "returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query."},
	}

	registered := Helpers().Get()
	for _, h := range registered {
		h.Category = "Extensions"
	}
	ret = append(ret, registered...)

	return
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
