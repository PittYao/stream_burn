package helper

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/duke-git/lancet/fileutil"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
)

func FixMu3u8(m3u8DirPath, m3u8FileName string) error {
	log.L.Sugar().Info("[开始修正m3u8]")
	m3u8FilePath := m3u8DirPath + "\\" + m3u8FileName
	exists := fileutil.IsExist(m3u8FilePath)
	if !exists {
		// 检测是否只有一个ts文件
		onlyOneTsFileInDir, tsFile, err := CheckOnlyOneTsFileInDir(m3u8DirPath)
		if err != nil {
			// 读取文件夹异常
			log.L.Sugar().Infof("[修正结束]:" + err.Error())
			return err
		}
		if !onlyOneTsFileInDir {
			// 多个videoxx.ts文件 则不处理
			return nil
		}

		// 添加默认m3u8文件
		if onlyOneTsFileInDir {
			// 文件夹下只有一个ts文件 新增一个初始化的m3u8文件
			log.L.Sugar().Infof("m3u8文件夹: %s 下只有一个ts文件", m3u8DirPath)
			AddDefaultM3u8File(m3u8FilePath, tsFile)
			return nil
		}
	}

	// 获取m3u8文件 最后一行数据
	lastLine, _, err := ReadLast2LineM3u8File(m3u8FilePath)
	if err != nil {
		log.L.Sugar().Infof(err.Error())
		return err
	}

	// 文件最后一行内容包含video的内容，则修正m3u8文件内容,修正规则在方法注释上
	if strings.Contains(lastLine, "video") {
		log.L.Sugar().Info(fmt.Sprintf("m3u8文件: %s 最后一行内容包含video,开始修正m3u8文件", m3u8FilePath))
		// 获取文件夹下最后的ts文件
		err, lastTsFile := GetLastModifyTsFileFromDir(m3u8DirPath)
		if err != nil {
			return err
		}

		AppendEndFlag2M3u8File(lastLine, m3u8FilePath, lastTsFile)
		log.L.Sugar().Infof("[修正结束]")
	}
	return nil
}

// ReadLast2LineM3u8File 获取m3u8文件 最后一行和倒数第二行数据
func ReadLast2LineM3u8File(filePath string) (string, string, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("打开m3u8文件异常: %s", filePath))
	}
	defer fi.Close()

	// 文件最后一行的数据
	var lastLine string
	var secondLastLine string

	br := bufio.NewReader(fi)
	for {
		a, c := br.ReadString('\n')
		if c == io.EOF {
			// 正常读取完成
			break
		}

		secondLastLine = lastLine
		lastLine = a
	}

	log.L.Sugar().Infof("m3u8: %s 最后一行内容:%s", filePath, lastLine)

	return lastLine, secondLastLine, nil
}

// AppendEndFlag2M3u8File m3u8文件最后添加结尾标志符
func AppendEndFlag2M3u8File(lastLine, m3u8FilePath string, lastTsFile os.FileInfo) {
	file, _ := os.OpenFile(m3u8FilePath, os.O_APPEND, 0777)
	defer file.Close()

	// 最后是最新的ts文件，只添加endlist
	if lastLine == lastTsFile.Name() {
		file.WriteString("#EXT-X-ENDLIST\n")
		log.L.Sugar().Infof(fmt.Sprintf("修正m3u8文件: %s,加入结束标志符", m3u8FilePath))
		return
	}

	// 最后不是最新的ts文件，添加最后一个文件的信息

	// 获取最新ts的文件时长
	time, err := GetFileModifyTime(lastTsFile)
	if err != nil || time == "0" {
		// 异常默认60秒时长
		file.WriteString("#EXTINF:60.000000,\n")
		log.L.Sugar().Infof(fmt.Sprintf("修正m3u8文件: %s,加入默认60秒文件体", m3u8FilePath))
	} else {
		// 真实文件时长
		tsFileInfo := fmt.Sprintf("#EXTINF:" + time + ".000000,\n")
		file.WriteString(tsFileInfo)
		log.L.Sugar().Infof(fmt.Sprintf("修正m3u8文件: %s,加入最新文件信息: %s", m3u8FilePath, tsFileInfo))
	}
	file.WriteString(lastTsFile.Name() + "\n")
	file.WriteString("#EXT-X-ENDLIST\n")

}

// GetFileModifyTime 获取文件更新时长
func GetFileModifyTime(file os.FileInfo) (string, error) {

	// TODO linux环境下代码如下
	//linuxFileAttr := fileInfo.Sys().(*syscall.Stat_t)
	//createTime := linuxFileAttr.Ctim.Sec
	//endTime := linuxFileAttr.Mtim.Sec
	//duration := endTime.Sub(create)

	// windows下代码如下
	winFileAttr := file.Sys().(*syscall.Win32FileAttributeData)
	createTime := time.Unix(winFileAttr.CreationTime.Nanoseconds()/1e9, 0)
	endTime := time.Unix(winFileAttr.LastWriteTime.Nanoseconds()/1e9, 0)
	duration := endTime.Sub(createTime)

	return fmt.Sprint(duration.Seconds()), nil
}

// GetLastModifyTsFileFromDir 获取m3u8文件夹下最后更新的一个ts文件
func GetLastModifyTsFileFromDir(m3u8Dir string) (error, os.FileInfo) {
	// 1.读取文件夹下所有文件且根据文件修改时间排序，并排除m3u8文件
	sortFiles, err := ReadDirSortModTimeExcludeM3u8(m3u8Dir)
	// 排序情况
	// video001.ts---2022-04-19 09:13:20.9309853 +0800 CST
	// video000.ts---2022-04-19 09:12:20.9653066 +0800 CST

	if err != nil {
		log.L.Sugar().Infof(err.Error())
		return err, nil
	}

	// 返回最新修改的文件
	return err, sortFiles[0]
}

// ReadDirSortModTimeExcludeM3u8 读取文件夹下所有文件且根据文件修改时间排序，并排除m3u8文件
func ReadDirSortModTimeExcludeM3u8(dirname string) ([]fs.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("打开m3u8文件夹:%s 异常", dirname))
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("遍历m3u8文件夹:%s 异常", dirname))
	}

	if len(list) == 0 {
		return nil, errors.New(fmt.Sprintf("m3u8文件夹:%s 下没有任何文件存在", dirname))
	}

	var subList []os.FileInfo
	for _, fileInfo := range list {
		if !strings.Contains(fileInfo.Name(), "playlist.m3u8") {
			subList = append(subList, fileInfo)
		}
	}

	if len(subList) == 0 {
		return nil, errors.New(fmt.Sprintf("m3u8文件夹:%s 下没有ts文件存在", dirname))
	}

	list = SortFileByModifyTime(subList)
	return list, nil
}

// SortFileByModifyTime 根据文件夹中文件修改时间排降序
func SortFileByModifyTime(pl []os.FileInfo) []os.FileInfo {
	sort.Slice(pl, func(i, j int) bool {
		flag := false
		if pl[i].ModTime().After(pl[j].ModTime()) {
			flag = true
		} else if pl[i].ModTime().Equal(pl[j].ModTime()) {
			if pl[i].Name() < pl[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return pl
}

// CheckOnlyOneTsFileInDir 检测文件夹下是否只有一个video000.ts,没有m3u8文件
func CheckOnlyOneTsFileInDir(dirPath string) (bool, os.FileInfo, error) {
	log.L.Sugar().Info("m3u8文件不存在,检测是否只有一个ts文件")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return false, nil, errors.New(fmt.Sprintf("m3u8文件不存在,读取m3u8文件夹: %s 目录下文件异常: %s", dirPath, err.Error()))
	}

	if len(files) == 0 {
		return false, nil, errors.New(fmt.Sprintf("m3u8文件不存在,读取m3u8文件夹: %s 目录下没有任何文件存在", dirPath))
	}

	if len(files) == 1 && files[0].Name() == "video000.ts" {
		return true, files[0], nil
	}

	log.L.Sugar().Info(fmt.Sprintf("m3u8文件不存在,m3u8文件夹: %s目录下有多个videoxx.ts文件", dirPath))
	return false, nil, nil
}

// AddDefaultM3u8File 拷贝m3u8文件
func AddDefaultM3u8File(distM3u8FilePath string, tsFile os.FileInfo) error {
	srcFilePath := "lib/m3u8/" + "playlist.m3u8"

	// 拷贝初始m3u8文件到指定路径
	if err := CopyFile(distM3u8FilePath, srcFilePath); err != nil {
		log.L.Sugar().Infof("拷贝初始化m3u8文件失败: %s", err.Error())
		return err
	}

	AppendEndFlag2M3u8File("", distM3u8FilePath, tsFile)
	log.L.Sugar().Infof("[修正结束] m3u8文件:%s下只有一个ts文件,新增了m3u8的初始文件: %s", distM3u8FilePath, tsFile.Name())

	return nil
}
