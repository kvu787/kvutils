package kvutils

import (
	"time"
)

func BusySleep(duration time.Duration) {
	end := time.Now().Add(duration)
	for time.Now().Before(end) {
	}
}
