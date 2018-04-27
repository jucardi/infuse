package gotmpl

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/jucardi/go-infuse/util/loader"
	"github.com/jucardi/go-infuse/util/log"
	"github.com/jucardi/go-strings"
)

// GoTemplate encapsulates the template context
type GoTemplate struct {
	definitions map[string]string
	template    string
	name        string
}

// Name represents the name of the ITemplate instance. This name will be used internally when creating the go template,
// so it can be used as a reference from other definitions when using the {{template}} directive.
func (t *GoTemplate) Name() string {
	return t.name
}

// Type returns the template type of this instance.
func (t *GoTemplate) Type() string {
	return "go"
}

// ParseJSON parses the template using the string representation of a JSON
func (t *GoTemplate) ParseJSON(writer io.Writer, jsonStr string) error {
	val, err := loader.LoadJSON(jsonStr)
	if err != nil {
		return fmt.Errorf("an error ocurred while unmarshalling the JSON, %v", err)
	}
	return t.Parse(writer, val)
}

// Parse parses the template
func (t *GoTemplate) Parse(writer io.Writer, data interface{}) error {
	log.Debug(" <-- templates.Template.Parse")
	str := t.prepare()
	log.Debug(" === template string ===\n\n", str, " =======================")
	tmpl := template.New(t.name).Funcs(defaultFuncMap())
	if _, err := tmpl.Parse(str); err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

// LoadFileTemplate loads the given file as the template to be parsed.
func (t *GoTemplate) LoadFileTemplate(filename string) error {
	tmplStr, err := loader.LoadTemplate(filename)
	if err != nil {
		return fmt.Errorf("unable to load file '%s', %v", filename, err)
	}
	return t.LoadTemplate(tmplStr)
}

// LoadTemplate loads the given string as the template to be parsed.
func (t *GoTemplate) LoadTemplate(tmpl string) error {
	return validate(t.name, tmpl, func() {
		t.template = tmpl
	})
}

// LoadFileDefinitionsByPattern uses a pattern to find the file definitions to be loaded for the template parsing
func (t *GoTemplate) LoadFileDefinitionsByPattern(pattern string) error {
	result, err := loader.LoadTemplates(pattern)
	if err != nil {
		return err
	}
	for k, v := range result {
		t.definitions[k] = v
	}
	return nil
}

// LoadFileDefinition loads a file(s) as definition(s) {{define "filename"}} using the filename as the name for the definition, to be used for 'template' directives
func (t *GoTemplate) LoadFileDefinition(files ...string) error {
	for _, filepath := range files {
		split := strings.Split(filepath, "/")
		split = strings.Split(split[len(split)-1], "\\")
		filename := split[len(split)-1]

		if tmplStr, err := loader.LoadTemplate(filepath); err != nil {
			return err
		} else if err := t.LoadDefinition(filename, tmplStr); err != nil {
			return err
		}
	}
	return nil
}

// LoadDefinition loads the give template string as a definition {{define "name"}}, using the given name as the name of the definition, to be used for 'template' directives.
func (t *GoTemplate) LoadDefinition(name, tmpl string) error {
	return validate(name, tmpl, func() {
		t.definitions[name] = tmpl
	})
}

func (t *GoTemplate) prepare() string {
	builder := stringx.Builder()

	for k, v := range t.definitions {
		if k == t.name {
			continue
		}
		builder.
			AppendLinef("{{define \"%s\"}}", k).
			AppendLine(v).
			AppendLine("{{end}}")
	}
	return builder.AppendLine(t.template).Build()
}

// New creates a new template utility which extends the default built in functions for Go templates.
func New(name ...string) *GoTemplate {
	return &GoTemplate{
		definitions: map[string]string{},
		name:        stringx.GetOrDefault("base", name...),
	}
}

func validate(name, tmpl string, successFn func()) error {
	_, err := template.New(name).Funcs(defaultFuncMap()).Parse(tmpl)

	if err != nil {
		return fmt.Errorf("unable to load definition '%s', %v", name, err)
	}

	successFn()
	return nil
}
