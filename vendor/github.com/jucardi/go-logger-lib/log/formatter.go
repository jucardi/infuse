package log

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

const (
	TemplateDefault = `{{ if .LoggerName }}{{ MatchSize .LoggerName 10 }} | {{ end }}{{ .Level }} | {{ TimeFormat .Timestamp "HH:mm:ss" }} | {{ .Message }}`
)

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.Data["msg"]`. The message passed from Info, Warn, Error ..
// * `entry.Data["time"]`. The timestamp.
// * `entry.Data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type IFormatter interface {
	Format(io.Writer, *Entry) error
	SetTemplate(string) error
}

// BaseFormatter base structure for formatters
type BaseFormatter struct {
	templateHandler *template.Template
	helpers         template.FuncMap
}

func (f *BaseFormatter) SetTemplate(tmpl string) error {
	t, err := template.New("formatter").Funcs(f.helpers).Parse(tmpl)
	if err != nil {
		f.SetTemplate(TemplateDefault)
		fmt.Fprintf(os.Stderr, "error occurred while setting template for logger, %s  > ", err.Error())
		fmt.Fprintln(os.Stderr, "setting default template")
		return err
	}
	f.templateHandler = t
	return nil
}
