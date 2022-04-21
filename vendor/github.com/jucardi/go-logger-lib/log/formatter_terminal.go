package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// BaseTerminalFormatter base structure to create formatters for a terminal
type BaseTerminalFormatter struct {
	BaseFormatter
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool
	supportsColor             *bool
	theme                     *TerminalTheme
}

func (f *BaseTerminalFormatter) isColored() bool {
	if f.supportsColor == nil {
		supportsColor := f.ForceColors

		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
			supportsColor = true
		} else if ok && force == "0" {
			supportsColor = false
		} else if os.Getenv("CLICOLOR") == "0" {
			supportsColor = false
		} else if strings.Contains(os.Getenv("TERM"), "color") {
			supportsColor = true
		}
		f.supportsColor = &supportsColor
	}

	return *f.supportsColor && !f.DisableColors
}

func (f *BaseTerminalFormatter) SetTheme(scheme *TerminalTheme) {
	f.theme = scheme
	if scheme.Template != "" {
		f.SetTemplate(scheme.Template)
	}
}

// TextFormatter formats logs into text
type TerminalFormatter struct {
	BaseTerminalFormatter
}

func NewTerminalFormatter() *TerminalFormatter {
	ret := &TerminalFormatter{}
	ret.helpers = getDefaultHelpers()
	ret.SetTheme(TerminalThemeDefault)
	return ret
}

// Format renders a single log entry
func (f *TerminalFormatter) Format(writer io.Writer, entry *Entry) error {
	if f.templateHandler == nil {
		return errors.New("no template parser found")
	}

	if writer == nil {
		return errors.New("writer cannot be nil")
	}

	entry.AddMetadata(metadataColorEnabled, f.isColored())
	entry.AddMetadata(metadataColorScheme, f.theme.Schemes)
	if err := f.templateHandler.Execute(writer, entry); err != nil {
		return fmt.Errorf("unable to write log to io writer, %s", err.Error())
	}

	fmt.Fprintln(writer)
	return nil
}
