package util

import "time"

// Makes times compatible with old Saerro API
func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
