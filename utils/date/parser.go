package date

import (
	"time"
)

func DefaultToTime(value string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	return t, err
}

func DefaultToTimeOrDefault(value string) time.Time {
	t, err := DefaultToTime(value)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}
