package usecase

import "time"

func defaultNowUTC() time.Time {
	return time.Now().UTC()
}
