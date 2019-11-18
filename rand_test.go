package util

import (
	"testing"
)

func TestGetRandUint64(t *testing.T) {
	for i := 0; i < 5; i++ {
		randomUint64, err := GetRandUint64()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(randomUint64)
	}
}
