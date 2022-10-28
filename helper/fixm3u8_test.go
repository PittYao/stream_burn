package helper

import (
	"os"
	"testing"
)

func TestGetFileModifyTime(t *testing.T) {
	fileInfo, _ := os.Stat("D:\\videodata\\192.168.99.215-192.168.99.215-192.168.99.215\\2022.10.17-08.59.56\\video048.ts")
	time, _ := GetFileModifyTime(fileInfo)
	println(time)
}
