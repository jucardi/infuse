package gotmpl

import (
	"fmt"
	"io"
	"text/template"

	"github.com/jucardi/go-strings"
	"github.com/jucardi/go-infuse/util/log"
	"github.com/jucardi/go-infuse/util/loader"
)

// Template encapsulates the template context
type Template struct {
	definitions map[string]string
	template    string
	name        string
}

// ParseJSON parses the template using the string representation of a JSON
func (t *Template) ParseJSON(writter io.Writer, jsonStr string) error {
	val, err := loader.LoadJSON(jsonStr)
	if err != nil {
		return fmt.Errorf("an error ocurred while unmarshalling the JSON, %v", err)
	}
	return t.Parse(writter, val)
}

// Parse parses the template
func (t *Template) Parse(writer io.Writer, data interface{}) error {
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
func (t *Template) LoadFileTemplate(filename string) error {
	tmplStr, err := loader.LoadTemplate(filename)
	if err != nil {
		return fmt.Errorf("unable to load file '%s', %v", filename, err)
	}
	return t.LoadTemplate(tmplStr)
}

// LoadTemplate loads the given string as the template to be parsed.
func (t *Template) LoadTemplate(tmpl string) error {
	return validate(t.name, tmpl, func() {
		t.template = tmpl
	})
}

// LoadFileDefinitionsByPattern uses a pattern to find the file definitions to be loaded for the template parsing
func (t *Template) LoadFileDefinitionsByPattern(pattern string) error {
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
func (t *Template) LoadFileDefinition(files ...string) error {
	for _, filename := range files {
		if tmplStr, err := loader.LoadTemplate(filename); err != nil {
			return err
		} else if err := t.LoadDefinition(filename, tmplStr); err != nil {
			return err
		}
	}
	return nil
}

// LoadDefinition loads the give template string as a definition {{define "name"}}, using the given name as the name of the definition, to be used for 'template' directives.
func (t *Template) LoadDefinition(name, tmpl string) error {
	return validate(name, tmpl, func() {
		t.definitions[name] = tmpl
	})
}

func (t *Template) prepare() string {
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
func New(name string) *Template {
	return &Template{
		definitions: map[string]string{},
		name:        name,
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
