package tmplfunc

import (
	"time"
)

func NowFunc(now time.Time) func() time.Time {
	return func() time.Time { return now }
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func FormatTime(format string, t time.Time) string {
	return t.Format(format)
}

func TimeIn(tz string, t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return t, err
	}

	return t.In(loc), nil
}
