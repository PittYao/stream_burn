package helper

import (
	"testing"
)

func TestFreeDiskOverLimit(t *testing.T) {

}

func TestLoadDiskInfo(t *testing.T) {
	capacity, usedCapacity, freeCapacity, percent, err := LoadDiskInfo("f:")
	println(capacity, usedCapacity, freeCapacity, percent)
	println(err == nil)
}
