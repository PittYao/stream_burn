package helper

import (
	"errors"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/kirinlabs/HttpRequest"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
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

func CopyDir(src string, dest string, removeSrc bool) error {
	src = FormatPath(src)
	dest = FormatPath(dest)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("xcopy", src, dest, "/I", "/Y")
	case "darwin", "linux":
		cmd = exec.Command("cp", "-R", src, dest)
	}

	_, e := cmd.Output()
	if e != nil {
		log.L.Error("复制文件夹异常: " + e.Error())
		return errors.New("复制文件夹异常")
	}
	log.L.Info("复制文件成功", zap.String("src", src), zap.String("dest", dest))

	if removeSrc {
		// 复制完毕，删除源文件
		if err := os.RemoveAll(src); err != nil {
			log.L.Error("删除源文件失败: " + src)
		} else {
			log.L.Info("删除源文件成功: " + src)
		}
	}

	return nil
}

func FormatPath(s string) string {
	switch runtime.GOOS {
	case "windows":
		return strings.Replace(s, "/", "\\", -1)
	case "darwin", "linux":
		return strings.Replace(s, "\\", "/", -1)
	default:
		log.L.Info("only support linux,windows,darwin, but os is " + runtime.GOOS)
		return s
	}
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

func DownloadFile2Path(url, savePath string) (error, string) {
	request := HttpRequest.NewRequest()

	res, err := request.Get(url)
	statusCode := res.StatusCode()
	if statusCode == http.StatusNotFound || statusCode == http.StatusInternalServerError || err != nil {
		log.L.Error("下载文件失败,访问url失败")
		return errors.New("下载文件失败，访问url失败"), ""
	}

	body, err := res.Body()

	// url中获取文件名
	r, _ := http.NewRequest("GET", url, nil)
	fileName := path.Base(r.URL.Path)

	savePath = savePath + fileName

	err = os.WriteFile(savePath, body, os.ModePerm)
	if err != nil {
		log.L.Error("下载文件失败,存储到本地失败")
		return errors.New("下载文件失败,存储到本地失败"), ""
	}

	return nil, savePath

}
