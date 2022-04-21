package parser

import (
	"errors"
	"fmt"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/util/ioutils"
	"io"
	"io/ioutil"
	"os"
)

var (
	IgnoreList = []string{
		".DS_Store",
		".git",
		".idea",
		".vscode",
	}
)

// TemplateRequest encapsulates the required information to parse a template
type TemplateRequest struct {
	// DEPRECATED, Use `Path` instead
	Filename        string
	Path            string
	String          string
	Files           []string
	URL             string
	Output          string
	Definitions     []string
	SearchPattern   string
	ContinueOnError bool
}

func (t TemplateRequest) validate() error {
	if len(t.Files) == 0 && t.String != "" && t.URL != "" {
		return errors.New("only one input method allowed, specify either an input filename, a string or a URL")
	}
	if t.Path == "" {
		return errors.New("template path is required")
	}
	return nil
}

func (t TemplateRequest) load() (Data, error) {
	if t.Path == "" && t.Filename != "" {
		t.Path = t.Filename
	}

	if err := t.validate(); err != nil {
		return nil, err
	}

	dataObj := Data{}

	if len(t.Files) > 0 {
		for _, file := range t.Files {
			if err := dataObj.LoadFile(file); err != nil {
				return nil, err
			}
		}
	}

	if t.URL != "" {
		if err := dataObj.LoadURL(t.URL); err != nil {
			return nil, err
		}
	}
	return dataObj, nil
}

// Parse parses the given template with the given information
func Parse(req TemplateRequest) error {
	log.SetLevel(log.DebugLevel)
	data, err := req.load()

	if err != nil {
		return fmt.Errorf("unable to load data, %v", err)
	}

	stat, err := os.Stat(req.Path)

	if err != nil {
		return err
	}

	if stat.IsDir() {
		return parseDir(data, req)
	}

	return parseFile(data, req)
}

func readPath(path string, makeDir bool) (exists, isDir bool, err error) {
	println("reading path: ", path)
	if path == "" {
		return
	}

	_, err = os.Stat(path)

	if err != nil && os.IsNotExist(err) && makeDir {
		err = os.Mkdir(path, 0644)
	}
	if err != nil {
		return false, false, err
	}

	return true, true, nil
}

func parseDir(data Data, req TemplateRequest) error {
	if req.Output != "" {
		stat, err := os.Stat(req.Output)
		if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(req.Output, 0755)
		} else if err != nil {
			return err
		}

		if stat != nil && !stat.IsDir() {
			return errors.New("output must be a directory if input is also a directory")
		}
	}

	items, err := ioutil.ReadDir(req.Path)

	if err != nil {
		return err
	}

	for _, f := range items {
		if streams.From(IgnoreList).Contains(f.Name()) {
			continue
		}
		output := ""
		if req.Output != "" {
			output = paths.Combine(req.Output, f.Name())
		}
		newReq := TemplateRequest{
			Path:            paths.Combine(req.Path, f.Name()),
			String:          req.String,
			Files:           req.Files,
			URL:             req.URL,
			Output:          output,
			Definitions:     req.Definitions,
			SearchPattern:   req.SearchPattern,
			ContinueOnError: req.ContinueOnError,
		}

		var parser func(Data, TemplateRequest) error

		if f.IsDir() {
			parser = parseDir
		} else {
			parser = parseFile
		}

		if err := parser(data, newReq); err != nil {
			if req.ContinueOnError {
				log.Error(err)
			} else {
				return err
			}
		}
	}
	return nil
}

func parseFile(data Data, req TemplateRequest) error {
	template := templates.Factory().New(req.Path)
	var writer io.WriteCloser

	// Load template
	if err := template.LoadFileTemplate(req.Path); err != nil {
		return fmt.Errorf("failed to load template '%s', %v", req.Path, err)
	}

	// Load template definitions.
	if len(req.Definitions) > 0 {
		if err := template.LoadFileDefinition(req.Definitions...); err != nil {
			return fmt.Errorf("failed to load definitions, %v", err)
		}
	}

	// Load definitions by search pattern
	if req.SearchPattern != "" {
		if err := template.LoadFileDefinitionsByPattern(req.SearchPattern); err != nil {
			return fmt.Errorf("failed to load definitions, %v", err)
		}
	}

	// Establish the io.Writer to use
	if req.Output != "" {
		if w, err := ioutils.NewFileWriter(req.Output); err != nil {
			return fmt.Errorf("unable to open file '%s', %v", req.Output, err)
		} else {
			writer = w
			defer w.Close()
		}
	} else {
		writer = os.Stdout
	}

	if err := template.Parse(writer, data.ToMap()); err != nil {
		return fmt.Errorf("failed to parse the template, %v", err)
	}

	return nil
}
