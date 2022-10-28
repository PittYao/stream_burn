package helper

import (
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/duke-git/lancet/datetime"
	"github.com/duke-git/lancet/fileutil"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// EncodeRtspUrl 替换rtspUrl中的&为 %26
func EncodeRtspUrl(rtspUrl string) string {
	index := strings.Index(rtspUrl, "&")
	if index != -1 {
		rtspUrl = strings.Replace(rtspUrl, "&", "%26", -1)
	}
	return rtspUrl
}

// CmdString2Array cmd命令由string转为string数组
func CmdString2Array(strCmd string) []string {
	args := strings.Split(strCmd, " ")
	var resultArgs []string
	for _, v := range args {
		if v != "" {
			resultArgs = append(resultArgs, v)
		}
	}
	return resultArgs
}

// GetRebootTime 计算重启时间，用于定时重启
func GetRebootTime() time.Time {
	//minute := config.C.Video.M3u8MaxTime * 60 // 单位分钟
	return datetime.AddMinute(time.Now(), 1)
}

// GetRtmpUrl 由rtsp地址获取rtmp地址
func GetRtmpUrl(rtspUrl string) string {
	rtmpUrl := "rtmp://" + config.C.Server.Ip + ":1935/live/" + Md5ByString(rtspUrl)
	return rtmpUrl
}

// GetM3u8Url 获取m3u8网络地址
func GetM3u8Url(ip, filePath string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	index := strings.Index(filePath, "/")
	if index == -1 {
		log.L.Sugar().Errorf("m3u8Url生成异常 filePath:%s", filePath)
		return ""
	}

	m3u8Url := fmt.Sprintf("http://%s%s%s/%s", ip, consts.M3u8UrlPort, filePath[index:], consts.M3u8File)

	// http://192.168.99.19:8880/videodata/publicSingle/192.168.99.117/2022.04.21-13.08.08/playlist.m3u8
	return m3u8Url
}

// GetTempM3u8Url 获取m3u8临时网络地址
func GetTempM3u8Url(ip, filePath, m3u8FileName string) string {
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	index := strings.Index(filePath, "/")
	if index == -1 {
		log.L.Sugar().Errorf("m3u8Url生成异常 filePath:%s", filePath)
		return ""
	}

	m3u8Url := fmt.Sprintf("http://%s%s%s/%s", ip, consts.M3u8UrlPort, filePath[index:], m3u8FileName)

	// http://192.168.99.19:8880/videodata/publicSingle/192.168.99.117/2022.04.21-13.08.08/playlist.m3u8
	return m3u8Url
}

func M3u8TempUrlHandler(savePath string) (string, error) {
	srcM3u8Path := savePath + "\\" + consts.M3u8File
	if !fileutil.IsExist(srcM3u8Path) {
		log.L.Sugar().Errorf("路径下m3u8文件不存在:%s,id:%d", savePath)
		return "", errors.New("路径下m3u8文件不存在")
	}

	// 复制m3u8文件
	m3u8FileName := RandomStr() + consts.M3u8File
	distM3u8Path := savePath + "\\" + m3u8FileName
	if err := CopyFile(distM3u8Path, srcM3u8Path); err != nil {
		log.L.Error("拷贝m3u8文件失败", zap.Any("err", err))
		return "", errors.New("拷贝m3u8文件失败")
	}
	// 修复m3u8临时文件
	FixMu3u8(savePath, m3u8FileName)
	// 构建m3u8网络地址
	m3u8Url := GetTempM3u8Url(config.C.Server.Ip, savePath, m3u8FileName)
	return m3u8Url, nil
}

// RedirectUrlBuilder 重定向url
func RedirectUrlBuilder(serverHost, port, url string) string {
	// http://192.168.99.127:8008/record/single/start
	return fmt.Sprintf("%s%s%s%s", consts.Http, serverHost, port, url)
}

// CalculatingTime 请求的开始结束时间 和 视频的开始结束时间 构建ffmpeg下载参数
func CalculatingTime(startTime, endTime, fileStartTime, fileEndTime *time.Time) (string, string) {
	// 1.转化为秒
	startTimeUnix := startTime.Unix()
	endTimeUnix := endTime.Unix()
	fileStartTimeUnix := fileStartTime.Unix()
	fileEndTimeUnix := fileEndTime.Unix()

	// 2.比较大小

	// a1 = fileStartTimeUnix ; a2 = fileEndTimeUnix
	// b1 = startTimeUnix ; b2 = endTimeUnix

	// b1 >= a1 && b2 <= a2  ===>  b1~b2
	var ss int64
	var duration int64
	if startTimeUnix >= fileStartTimeUnix && endTimeUnix <= fileEndTimeUnix {
		// 距离文件的起始时间多少秒: b1-a1
		ss = startTimeUnix - fileStartTimeUnix
		// 截取多少秒: b2-b1
		duration = endTimeUnix - startTimeUnix
	} else if startTimeUnix < fileStartTimeUnix && endTimeUnix > fileStartTimeUnix && endTimeUnix <= fileEndTimeUnix {
		// b1 < a1 && a1<=b2<=a2 ===> a1~b2
		// 距离文件的起始时间多少秒: 0
		ss = 0
		// 截取多少秒: b2-a1
		duration = endTimeUnix - fileStartTimeUnix
	} else if endTimeUnix > startTimeUnix && startTimeUnix >= fileStartTimeUnix && startTimeUnix < fileEndTimeUnix && endTimeUnix > fileEndTimeUnix {
		// b2 > b1 && a1<=b1<=a2 ===> b1~a2
		// 距离文件的起始时间多少秒: b1-a1
		ss = startTimeUnix - fileStartTimeUnix
		// 截取多少秒: a2-b1
		duration = fileEndTimeUnix - startTimeUnix
	} else if startTimeUnix < fileStartTimeUnix && fileEndTimeUnix < endTimeUnix {
		// b1 < a1 && a2 < b2
		// 距离文件的起始时间多少秒: a2-a1
		ss = 0
		// 截取多少秒: a2-b1
		duration = fileEndTimeUnix - fileStartTimeUnix
	}
	return strconv.FormatInt(ss, 10), strconv.FormatInt(duration, 10)
}

// GetSaveFileCmd 下载视频的ffmpeg命令构建
func GetSaveFileCmd(ss, m3u8Url, duration, mp4SavePath string) ([]string, string) {
	//-ss ${ss} -i ${m3u8Url} -c copy -t ${duration} -f segment -r 25 -segment_time 3600 -crf 28 ${mp4SavePath}
	cmdArgs := "  -ss " + ss + " -i " + m3u8Url + " -c copy -t " + duration + " -f segment -r 25 -segment_time 3600 -crf 28 " + mp4SavePath

	log.L.Info("构建下载视频命令",
		zap.Any("Ss", ss),
		zap.Any("M3u8Url", m3u8Url),
		zap.Any("Duration", duration),
		zap.Any("Mp4SavePath", mp4SavePath),
		zap.Any("ffmpeg_cmd", cmdArgs))

	return CmdString2Array(cmdArgs), cmdArgs
}
