package cis

import (
	"strconv"
	"strings"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		*t = Time{}
		return
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		*t = Time(time.UnixMilli(i).UTC())
		return
	}

	tp, err := time.Parse("Jan _2, 2006, 3:04:05 PM", s)

	*t = Time(tp.UTC())
	return
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
