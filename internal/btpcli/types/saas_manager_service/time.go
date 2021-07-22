package saas_manager_service

import (
	"strconv"
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
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.UnixMilli(i).UTC())

	return nil
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
