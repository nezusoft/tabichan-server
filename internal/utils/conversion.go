package utils

import "time"

func ConvertTimeStringToRFC3339(timeInput string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeInput)
	if err != nil {
		return time.Now(), err
	}
	return parsedTime, nil
}
