package service

import (
	"time"
)

func parsePeriodDates(period string) (time.Time, time.Time, int, int, error) {
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, err
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end, int(t.Month()), t.Year(), nil
}

func convertToFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	}
	return 0, false
}
