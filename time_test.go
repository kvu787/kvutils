package util

import (
	"testing"
	"time"
)

func TestBusySleep(t *testing.T) {
	BusySleep(3 * time.Second)
}
