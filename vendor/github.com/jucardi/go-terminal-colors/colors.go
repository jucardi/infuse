package fmtc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
)

type Color int

// special formats
const (
	Clear Color = iota
	Bold
	Dim
	Italic
	Underline
	Inverted
	Hidden
)

// foreground
const (
	Black Color = 30 + iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	Gray
)
const (
	DarkGray Color = 90 + iota
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	White
)

// background
const (
	BgBlack Color = 40 + iota
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgGray
)
const (
	BgDarkGray Color = 100 + iota
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgWhite
)

var colorMap = map[string]Color{
	"clear":     Clear,
	"bold":      Bold,
	"dim":       Dim,
	"italic":    Italic,
	"underline": Underline,
	"inverted":  Inverted,
	"hidden":    Hidden,

	"black":        Black,
	"red":          Red,
	"green":        Green,
	"yellow":       Yellow,
	"blue":         Blue,
	"magenta":      Magenta,
	"cyan":         Cyan,
	"gray":         Gray,
	"darkgray":     DarkGray,
	"lightred":     LightRed,
	"lightgreen":   LightGreen,
	"lightyellow":  LightYellow,
	"lightblue":    LightBlue,
	"lightmagenta": LightMagenta,
	"lightcyan":    LightCyan,
	"white":        White,

	"bgblack":        BgBlack,
	"bgred":          BgRed,
	"bggreen":        BgGreen,
	"bgyellow":       BgYellow,
	"bgblue":         BgBlue,
	"bgmagenta":      BgMagenta,
	"bgcyan":         BgCyan,
	"bggray":         BgGray,
	"bgdarkgray":     BgDarkGray,
	"bglightred":     BgLightRed,
	"bglightgreen":   BgLightGreen,
	"bglightyellow":  BgLightYellow,
	"bglightblue":    BgLightBlue,
	"bglightmagenta": BgLightMagenta,
	"bglightcyan":    BgLightCyan,
	"bgwhite":        BgWhite,
}

// Parse attempts to parse a string representing a predefined color.
func Parse(color string) (Color, error) {
	if v, ok := colorMap[stringx.New(color).ToLower().TrimSpace().S()]; ok {
		return v, nil
	}

	var l Color
	return l, fmt.Errorf("not a valid color: %q", color)
}

func doColors(colors []Color, w ...io.Writer) {
	joinedColors := ""
	if len(colors) > 0 {
		joinedColors = strings.Join(
			streams.From(colors).
				Map(colorToStr).
				ToArray().([]string),
			";",
		)
	}

	if len(w) > 0 && w[0] != nil {
		fmt.Fprint(w[0], esc, joinedColors, "m")
	} else {
		fmt.Print(esc, joinedColors, "m")
	}
}

func getColors(colors []Color) string {
	b := &bytes.Buffer{}
	doColors(colors, b)
	return b.String()
}

func wrapColors(colors []Color, str string) string {
	return getColors(colors) + str + clear
}

func colorToStr(i interface{}) interface{} {
	return strconv.Itoa(int(i.(Color)))
}
