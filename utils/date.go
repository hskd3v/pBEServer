package utils

import "time"

const (
	dateLayout = time.RFC3339
)

// GetNow gets now
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString gets now as string
func GetNowString() string {
	return GetNow().Format(dateLayout)
}
