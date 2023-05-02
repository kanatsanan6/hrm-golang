package utils

import "time"

func StringToDateTime(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
