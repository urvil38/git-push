package utils

import (
	"testing"
)

func TestGetCurrentWorkingDirName(t *testing.T) {
	wanted := "utils"
	str := GetCurrentWorkingDirName()
	if str != wanted {
		t.Fatalf("Wanted %s but Get %s", wanted, str)
	}
	if str == "" {
		t.Fatalf("Wanted %s but Get empty dir name", wanted)
	}
}
