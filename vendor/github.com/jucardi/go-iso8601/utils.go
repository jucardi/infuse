package iso8601

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jucardi/go-strings/stringx"
)

var (
	DefaultMonthsMap = MonthsEng
	MonthsEng        = map[time.Month]string{
		time.January:   "January",
		time.February:  "February",
		time.March:     "March",
		time.April:     "April",
		time.May:       "May",
		time.June:      "June",
		time.July:      "July",
		time.August:    "August",
		time.September: "September",
		time.October:   "October",
		time.November:  "November",
		time.December:  "December",
	}
	MonthsEsp = map[time.Month]string{
		time.January:   "Enero",
		time.February:  "Febrero",
		time.March:     "Marzo",
		time.April:     "Abril",
		time.May:       "Mayo",
		time.June:      "Junio",
		time.July:      "Julio",
		time.August:    "Agosto",
		time.September: "Septiembre",
		time.October:   "Octubre",
		time.November:  "Noviembre",
		time.December:  "Diciembre",
	}
)

func TimeToIsoUtc(timestamp time.Time) string {
	return TimeToString(timestamp.UTC(), "yyyy-MM-ddTHH:mm:ssZ")
}

func TimeToString(timestamp time.Time, format string, monthsMap ...map[time.Month]string) string {
	// Converts the value of the current Date object to its equivalent string representation using the specified format
	// and the formatting conventions following the ISO 8601 standard.
	//
	// Example: "yyyy-MM-dd HH:mm:ss"

	year, month, day := timestamp.Date()
	hour, min, sec := timestamp.Clock()

	parsedHour12 := hour
	if hour == 0 {
		parsedHour12 = 12
	} else if hour > 12 {
		parsedHour12 = hour - 12
	}

	yearStr := fmt.Sprintf("%04d", year)
	year2DigitsStr := yearStr[len(yearStr)-2:]
	year2Digits, _ := strconv.Atoi(year2DigitsStr)
	mMap := DefaultMonthsMap
	if len(monthsMap) > 0 && monthsMap[0] != nil {
		mMap = monthsMap[0]
	}
	ret := stringx.New(format).
		Replace("HH", fmt.Sprintf("%02d", hour), -1).
		Replace("H", strconv.Itoa(hour), -1).
		Replace("hh", fmt.Sprintf("%02d", parsedHour12), -1).
		Replace("h", strconv.Itoa(parsedHour12), -1).
		Replace("mm", fmt.Sprintf("%02d", min), -1).
		Replace("m", strconv.Itoa(min), -1).
		Replace("ss", fmt.Sprintf("%02d", sec), -1).
		Replace("s", strconv.Itoa(sec), -1).
		Replace("dd", fmt.Sprintf("%02d", day), -1).
		Replace("d", strconv.Itoa(day), -1).
		Replace("yyyy", yearStr, -1).
		Replace("yy", year2DigitsStr, -1).
		Replace("y", strconv.Itoa(year2Digits), -1).
		Replace("M", "Mx", -1).
		Replace("MxMxMxMx", GetMonthString(month, false, mMap), -1).
		Replace("MxMxMx", GetMonthString(month, true, mMap), -1).
		Replace("MxMx", fmt.Sprintf("%02d", month), -1).
		Replace("Mx", strconv.Itoa(int(month)), -1).
		Replace("tt",
			func() string {
				if hour >= 12 {
					return "PM"
				}
				return "AM"
			}(), -1).
		S()

	return ret
}

func GetMonthString(month time.Month, shortMode bool, monthsMap ...map[time.Month]string) string {
	m := DefaultMonthsMap
	if len(monthsMap) > 0 && monthsMap[0] != nil {
		m = monthsMap[0]
	}
	ret, ok := m[month]
	if !ok {
		return fmt.Sprintf("Invalid month or mapping not found (%d)", month)
	}
	if shortMode {
		return ret[:3]
	}
	return ret
}
