package helper

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGetFileNamePrefix(t *testing.T) {
	//filename := "D:/test/通话记录.xlsx"
	filename := "通话记录.xlsx"
	preFix := GetFileNamePrefix(filename)
	println(preFix)
}

func TestGetTsFileNumber(t *testing.T) {
	filename := "video000.ts"
	println(GetTsFileNumber(filename))
}

func TestGetTsFileModifyTime(t *testing.T) {
	fileInfo, _ := os.Stat("D:\\videodata\\192.168.99.215-192.168.99.215-192.168.99.215\\2022.10.17-13.28.34\\video001.ts")
	fileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)

	create := time.Unix(fileSys.CreationTime.Nanoseconds()/1e9, 0)
	end := time.Unix(fileSys.LastWriteTime.Nanoseconds()/1e9, 0)
	access := time.Unix(fileSys.LastAccessTime.Nanoseconds()/1e9, 0)
	fmt.Println(create)
	fmt.Println(end)
	fmt.Println(access)

	duration := end.Sub(create)
	fmt.Println(duration.Seconds())

	duration1 := access.Sub(create)
	fmt.Println(duration1.Seconds())
}
