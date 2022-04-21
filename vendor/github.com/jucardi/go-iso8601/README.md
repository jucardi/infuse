# go-iso8601
Allows parsing from general ISO8601 expressions, from the Date only representations to the Repetitions and Durations or periods.

https://en.wikipedia.org/wiki/ISO_8601

Does not support week dates yet.

#### Getting started

To keep up to date with the most recent version:

```bash
go get github.com/jucardi/go-iso8601
```

#### Usage

###### Parsing a string

To parse the representation of an ISO8601 string simply use the `Parse` function.

Given the following expression

```Go
exp := "R5/2008-03-01T13:00:00Z/P1Y2M10DT2H30M/2017-03-01T13:00:00Z"
```

Where:
- `R5` represents 5 repetitions
- `2008-03-01T13:00:00Z` represents a start date in UTC
- `P1Y2M10DT2H30M` represents an interval or duration
- `2017-03-01T13:00:00Z` represents an end date in UTC

```Go
result, err := iso8601.Parse(exp)
```

The `Parse` function will return a struct representing the ISO8601 expression, which will be equal to:

```Go
result := &IntervalDescriptor {
	Start:   startTime, // a time.Time struct obtained by parsing the start date string
	End:     endTime,   // a time.Time struct obtained by parsing the end date string
	Repeats: 5,         // Equal to what was indicated in the Rn portion of the string
	Period: &Period{    // Struct containing the values defined by the interval or duration portion of the string
		Years:   1,
		Months:  2,
		Days:    10,
		Hours:   2,
		Minutes: 30,
		Seconds: 0,
	}
}
```

###### Converting to a string

By having a defined `IntervalDescriptor` struct, simply by doing a `.ToString()`, will return the string representation of the struct in ISO8601 format.

#### The `Period` struct

Works similar to it's parent `IntervalDecriptor`, a `Period` may be created from scratch or may be obtained by using the `PeriodFromString` function and
passing a Period representation of the ISO8601.

The `Period` struct provides additional utility functions.

- `Normalize`: If a period has values that can be converted into a greater full unit, it will do so. For example, if a period has 100 seconds, by invoking
`Normalize`, `Seconds` will be set to 40 and the whole minute subtracted from the seconds amount will be added to the `Minutes` value.
- `ToDuration`: Converts the struct into a `time.Duration` representation to easily be used with the structs in the `time` package.
- `HasTime`: Indicates whether the period has any time values (Hours, Minutes or Seconds).