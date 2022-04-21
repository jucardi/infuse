package iso8601

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"time"
	"unicode"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
)

const (
	valueNotFound        = "attempting to assign %s but no value found"
	periodTemplateString = `P{{if .Years}}{{.Years}}Y{{end}}{{if .Months}}{{.Months}}M{{end}}{{if .Weeks}}{{.Weeks}}W{{end}}{{if .Days}}{{.Days}}D{{end}}{{if .HasTime}}T{{end }}{{if .Hours}}{{.Hours}}H{{end}}{{if .Minutes}}{{.Minutes}}M{{end}}{{if .Seconds}}{{.Seconds}}S{{end}}`
)

var (
	values            = []rune{'Y', 'M', 'W', 'D', 'H', 'M', 'S'}
	periodTemplate, _ = template.New("period").Parse(periodTemplateString)
)

// Period represents the interval structure defined in the ISO8601
type Period struct {
	Years   int `json:"years" bson:"years"`
	Months  int `json:"months" bson:"months"`
	Weeks   int `json:"weeks" bson:"weeks"`
	Days    int `json:"days" bson:"days"`
	Hours   int `json:"hours" bson:"hours"`
	Minutes int `json:"minutes" bson:"minutes"`
	Seconds int `json:"seconds" bson:"seconds"`
}

// PeriodFromString creates a *Period by parsing the ISO8601 representation of a period.
func PeriodFromString(value string) (*Period, error) {
	runes := []rune(value)
	builder := stringx.Builder()
	timeEnabled := false
	ret := &Period{}
	for i, v := range runes {
		if i == 0 && v != 'P' {
			return nil, errors.New("invalid period representation, must start with P")
		}

		if v == 'P' {
			continue
		}

		if unicode.IsDigit(v) {
			if i == len(runes)-1 {
				return nil, errors.New("the last character cannot be a number")
			}
			builder.AppendRune(v)
			continue
		}

		if v == 'T' {
			timeEnabled = true
			continue
		}

		if !streams.FromArray(values).Contains(v) {
			return nil, fmt.Errorf("invalid value found, %s", strconv.QuoteRune(v))
		}

		if builder.IsEmpty() {
			return nil, fmt.Errorf(valueNotFound, strconv.QuoteRune(v))
		}
		val, _ := strconv.Atoi(builder.Build())
		builder = stringx.Builder()

		switch v {
		case 'Y':
			ret.Years = val
		case 'M':
			if timeEnabled {
				ret.Minutes = val
			} else {
				ret.Months = val
			}
		case 'W':
			ret.Weeks = val
		case 'D':
			ret.Days = val
		case 'H':
			if !timeEnabled {
				return nil, fmt.Errorf("found time component without time enabler 'T', %s", strconv.QuoteRune(v))
			}
			ret.Hours = val
		case 'S':
			if !timeEnabled {
				return nil, fmt.Errorf("found time component without time enabler 'T', %s", strconv.QuoteRune(v))
			}
			ret.Seconds = val
		}
	}

	return ret, nil
}

// Normalize normalizes the period values.
func (p *Period) Normalize() *Period {
	seconds := p.Seconds % 60
	minutes := p.Seconds/60 + p.Minutes
	minutes, hours := minutes%60, minutes/60+p.Hours
	hours = hours % 24
	d := hours/24 + p.Days + p.Weeks*7
	days := d % 30
	months, years := p.Months%12+d/30, p.Months/12+p.Years

	return &Period{years, months, 0, days, hours, minutes, seconds}
}

// ToDuration converts the period into a representation of `time.Duration`. Since periods do not have an actual date, Years are assumed to have 365 days and months
// are assumed to have 30 days.
func (p *Period) ToDuration() time.Duration {
	days := p.Years*365 + p.Months*30 + p.Days
	return time.Duration(days*24+p.Hours)*time.Hour + time.Duration(p.Minutes)*time.Minute + time.Duration(p.Seconds)*time.Second
}

// HasTime indicates whether the `Period` has a time component.
func (p *Period) HasTime() bool {
	return p.Hours > 0 || p.Minutes > 0 || p.Seconds > 0
}

// ToString converts the `Period` representation to a ISO8601 string,
func (p *Period) ToString() string {
	buf := new(bytes.Buffer)
	periodTemplate.Execute(buf, p)
	return buf.String()
}

// Apply period to timestamp and return result.
func (p *Period) Apply(t time.Time) time.Time {
	y, m, d, h, M, s, n := t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
	duration := time.Hour*time.Duration(p.Hours) + time.Minute*time.Duration(p.Minutes) + time.Second*time.Duration(p.Seconds)
	return time.Date(p.Years+y, time.Month(p.Months)+m, p.Days+d, h, M, s, n, t.Location()).Add(duration)
}
