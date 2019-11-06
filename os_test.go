package kvutils

import (
	"testing"
)

func TestDoesFileExist(t *testing.T) {
	doesExist, err := DoesFileExist("os.go")
	if !doesExist || err != nil {
		t.FailNow()
	}

	doesExist, err = DoesFileExist("nofile.go")
	if doesExist || err != nil {
		t.FailNow()
	}
}
