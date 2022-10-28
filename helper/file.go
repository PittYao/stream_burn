package helper

import (
	"github.com/PittYao/stream_burn/internal/consts"
	"os"
	"path"
)

// CopyFile 拷贝文件
func CopyFile(distFilePath string, srcFilePath string) error {
	bytesRead, err := os.ReadFile(srcFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(distFilePath, bytesRead, 0755)
	if err != nil {
		return err
	}
	return nil
}

// GetFileName 获取文件名称
func GetFileName(filePath string) string {
	return path.Base(filePath)
}

// GetFileNamePrefix 获取文件名称前缀
func GetFileNamePrefix(filePath string) string {
	fileName := GetFileName(filePath)
	suffix := GetFileNameSuffix(filePath)
	return fileName[0 : len(fileName)-len(suffix)]
}

// GetFileNameSuffix 获取文件名称后缀
func GetFileNameSuffix(filePath string) string {
	return path.Ext(filePath)
}

// GetTsFileNumber 获取ts文件数 video102.ts -> 102 video002.ts -> 2 video112.ts -> 112
func GetTsFileNumber(filePath string) string {
	filePrefix := GetFileNamePrefix(filePath)
	s := filePrefix[len(consts.TsFilePrefix):]

	// 去除0
	var a string
	for i := 0; i < len(s); i++ {
		item := s[i]

		if i == 0 {
			if string(item) == "0" {
				a = s[i+1:]
			} else {
				a = s
				break
			}
		}
		if i == 1 && string(item) == "0" {
			if string(s[0]) == "0" {
				a = s[i+1:]
			} else {
				break
			}
		}
	}
	return a
}
