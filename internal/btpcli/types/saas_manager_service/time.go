package saas_manager_service

import (
	"strings"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		*t = Time{}
		return nil
	}

	layout := "Jan 2, 2006, 3:04:05 PM"

	// Normalize narrow no-break space (U+202F) to a regular space.
	// Some locales use it as the AM/PM separator instead of a regular space.
	s = strings.ReplaceAll(s, "\u202f", " ")

	timeString, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	*t = Time(timeString.UTC())

	return nil
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
