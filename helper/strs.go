package helper

import (
	"github.com/duke-git/lancet/random"
	"regexp"
	"unsafe"
)

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// CompressStr 利用正则表达式压缩字符串，去除空格或制表符
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}
func RandomStr() string {
	uuid, err := random.UUIdV4()
	if err != nil {
		return ""
	}
	return uuid
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
