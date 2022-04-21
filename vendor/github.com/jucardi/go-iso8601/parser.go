package iso8601

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const iso8601TemplateString = "{{if .Repeats}}R{{.Repeats}}/{{end}}{{if .Start}}{{.GetStartString}}/{{end}}{{if .Period}}{{.GetPeriodString}}/{{end}}{{if .End}}{{.GetEndString}}{{end}}"

var iso8601Template, _ = template.New("iso8601").Parse(iso8601TemplateString)

// IntervalDescriptor represents the information held by an ISO8601 expression
type IntervalDescriptor struct {
	Start   time.Time `json:"start" bson:"start"`
	End     time.Time `json:"end" bson:"end"`
	Repeats int       `json:"repeats" bson:"repeats"`
	Period  *Period   `json:"period" bson:"period"`
}

// Parse Parses an ISO8601 expression
func Parse(expression string) (*IntervalDescriptor, error) {
	split := strings.Split(expression, "/")
	endSet := false
	repeatsSet := false
	durationSet := false
	ret := &IntervalDescriptor{}

	for i, v := range split {
		if strings.HasPrefix(v, "R") {
			if i != 0 {
				return nil, errors.New("repetitions component must be at the beginning of the string")
			}

			if len(v) == 1 && i == 0 {
				ret.Repeats = -1
			} else {
				r, err := strconv.Atoi(v[1:])

				if  err != nil {
					return nil, fmt.Errorf("unable to parse repetitions, %s", err.Error())
				}

				if r <= 0 {
					return nil, errors.New("repeat value must be greater than zero")
				}
				ret.Repeats = r
			}
			repeatsSet = true
			continue
		}

		if strings.HasPrefix(v, "P") {
			if durationSet {
				return nil, errors.New("invalid iso8601, more than one period component detected")
			}

			p, err := PeriodFromString(v)
			if err != nil {
				return nil, fmt.Errorf("invalid period, unable to parse, %s", err.Error())
			}
			ret.Period = p
			durationSet = true
			continue
		}

		t, err := time.Parse(time.RFC3339, v)

		if err != nil {
			return nil, fmt.Errorf("unable to parse time component, %s", err.Error())
		}

		if i == 0 || i == 1 && repeatsSet {
			ret.Start = t
			continue
		} else if endSet {
			return nil, errors.New("invalid iso8601, more than one end date detected")
		}

		ret.End = t
		endSet = true
	}

	return ret, nil
}

// ToString returns a string representation of the interval descriptor
func (i *IntervalDescriptor) ToString() string {
	buf := new(bytes.Buffer)
	iso8601Template.Execute(buf, i)
	str := buf.String()

	if strings.HasPrefix(str, "/") {
		str = str[1:]
	}

	if strings.HasSuffix(str, "/") {
		str = str[:len(str)-1]
	}

	return strings.Replace(str, "//", "/", -1)
}

// GetStartString returns a string representation of the start timestamp
func (i *IntervalDescriptor) GetStartString() string {
	zero := time.Time{}
	if i.Start == zero {
		return ""
	}
	return i.Start.Format(time.RFC3339)
}

// GetEndString returns a string representation of the end timestamp
func (i *IntervalDescriptor) GetEndString() string {
	zero := time.Time{}
	if i.End == zero {
		return ""
	}
	return i.End.Format(time.RFC3339)
}

// GetPeriodString returns a string representation of the period descriptor
func (i *IntervalDescriptor) GetPeriodString() string {
	return i.Period.ToString()
}
