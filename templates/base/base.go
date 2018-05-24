package base

import (
	"fmt"
	"io"
	"strings"

	"github.com/jucardi/infuse/util/loader"
)

// AbstractTemplate encapsulates the common functionality of an ITemplate implementation
type AbstractTemplate struct {
	IAbstractTemplateMembers
	Definitions map[string]string
	Template    string
	NameStr     string
}

// Name represents the name of the ITemplate instance. This name will be used internally when creating the go template,
// so it can be used as a reference from other definitions when using the {{template}} directive.
func (t *AbstractTemplate) Name() string {
	return t.NameStr
}

// ParseJSON parses the template using the string representation of a JSON
func (t *AbstractTemplate) ParseJSON(writer io.Writer, jsonStr string) error {
	val, err := loader.LoadJSON(jsonStr)
	if err != nil {
		return fmt.Errorf("an error ocurred while unmarshalling the JSON, %v", err)
	}
	return t.Parse(writer, val)
}

// LoadFileTemplate loads the given file as the template to be parsed.
func (t *AbstractTemplate) LoadFileTemplate(filename string) error {
	tmplStr, err := loader.LoadTemplate(filename)
	if err != nil {
		return fmt.Errorf("unable to load file '%s', %v", filename, err)
	}
	return t.LoadTemplate(tmplStr)
}

// LoadFileDefinitionsByPattern uses a pattern to find the file definitions to be loaded for the template parsing
func (t *AbstractTemplate) LoadFileDefinitionsByPattern(pattern string) error {
	result, err := loader.LoadTemplates(pattern)
	if err != nil {
		return err
	}
	for k, v := range result {
		t.Definitions[k] = v
	}
	return nil
}

// LoadFileDefinition loads a file(s) as definition(s) {{define "filename"}} using the filename as the name for the definition, to be used for 'template' directives
func (t *AbstractTemplate) LoadFileDefinition(files ...string) error {
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
