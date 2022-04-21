package log

import "github.com/jucardi/go-terminal-colors"

// LevelColorScheme represents the terminal colors associated to the level parsing for each logging level.
type LevelColorScheme map[Level][]fmtc.Color

// TerminalColorScheme defines terminal colors that are tied to a log level and a field.
type TerminalColorScheme map[string]LevelColorScheme

// TerminalTheme contains the logging theme configuration for terminal logging
type TerminalTheme struct {
	Template string
	Schemes  TerminalColorScheme
}

var (
	TerminalThemeDefault = &TerminalTheme{
		Template: `{{ LoggerName . }}{{ Level . }}{{ Timestamp " HH:mm:ss " . }} {{ .Message }}`,
		Schemes: TerminalColorScheme{
			"loggerName": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				InfoLevel:  []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				WarnLevel:  []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				FatalLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				PanicLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
			},
			"level": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Bold, fmtc.DarkGray},
				InfoLevel:  []fmtc.Color{fmtc.Bold, fmtc.Cyan},
				WarnLevel:  []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.Bold, fmtc.Red},
				FatalLevel: []fmtc.Color{fmtc.Bold, fmtc.Red},
				PanicLevel: []fmtc.Color{fmtc.Bold, fmtc.Red},
			},
			"timestamp": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.DarkGray},
				InfoLevel:  []fmtc.Color{fmtc.Cyan},
				WarnLevel:  []fmtc.Color{fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.Red},
				FatalLevel: []fmtc.Color{fmtc.Red},
				PanicLevel: []fmtc.Color{fmtc.Red},
			},
		},
	}

	TerminalThemeAlternative = &TerminalTheme{
		Template: `{{ LoggerName . }}{{ Scheme "level" (string " " .Level " ") . }}{{ Timestamp " HH:mm:ss " . }} {{ .Message }}`,
		Schemes: TerminalColorScheme{
			"loggerName": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				InfoLevel:  []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				WarnLevel:  []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				FatalLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
				PanicLevel: []fmtc.Color{fmtc.Bold, fmtc.Yellow},
			},
			"level": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Bold, fmtc.DarkGray},
				InfoLevel:  []fmtc.Color{fmtc.Bold, fmtc.White, fmtc.BgBlue},
				WarnLevel:  []fmtc.Color{fmtc.Black, fmtc.BgYellow},
				ErrorLevel: []fmtc.Color{fmtc.Bold, fmtc.White, fmtc.BgRed},
				FatalLevel: []fmtc.Color{fmtc.Bold, fmtc.White, fmtc.BgRed},
				PanicLevel: []fmtc.Color{fmtc.Bold, fmtc.White, fmtc.BgRed},
			},
			"timestamp": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.BgBlack, fmtc.DarkGray},
				InfoLevel:  []fmtc.Color{fmtc.BgBlack, fmtc.Cyan},
				WarnLevel:  []fmtc.Color{fmtc.BgBlack, fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.BgBlack, fmtc.Red},
				FatalLevel: []fmtc.Color{fmtc.BgBlack, fmtc.Red},
				PanicLevel: []fmtc.Color{fmtc.BgBlack, fmtc.Red},
			},
		},
	}

	TerminalThemeCliApp = &TerminalTheme{
		Template: `{{ Timestamp " HH:mm:ss " . }} {{ Message . "           " }}`,
		Schemes: TerminalColorScheme{
			"timestamp": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Gray},
				InfoLevel:  []fmtc.Color{fmtc.Cyan},
				WarnLevel:  []fmtc.Color{fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.Red},
				FatalLevel: []fmtc.Color{fmtc.Red},
				PanicLevel: []fmtc.Color{fmtc.Red},
			},
			"message": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Gray},
				InfoLevel:  []fmtc.Color{fmtc.White},
				WarnLevel:  []fmtc.Color{fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.LightRed},
				FatalLevel: []fmtc.Color{fmtc.LightRed},
				PanicLevel: []fmtc.Color{fmtc.LightRed},
			},
		},
	}

	TerminalThemeCliAppNoTime = &TerminalTheme{
		Template: `{{ Message . "           " }}`,
		Schemes: TerminalColorScheme{
			"message": LevelColorScheme{
				DebugLevel: []fmtc.Color{fmtc.Gray},
				InfoLevel:  []fmtc.Color{fmtc.White},
				WarnLevel:  []fmtc.Color{fmtc.Yellow},
				ErrorLevel: []fmtc.Color{fmtc.LightRed},
				FatalLevel: []fmtc.Color{fmtc.LightRed},
				PanicLevel: []fmtc.Color{fmtc.LightRed},
			},
		},
	}
)
