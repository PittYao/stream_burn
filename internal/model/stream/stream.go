package stream

import (
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/model/burninfocmd"
	"github.com/PittYao/stream_burn/internal/model/burnsetting"
	"path"
	"strconv"
	"strings"
	"time"
)

type ToTask interface {
	ModelToTask() *Task
}

type Task struct {
	M3u8Url             string
	FfmpegSaveStartTime *time.Time // 存流开始时间
	FfmpegSaveCloseTime *time.Time // 存流结束时间
}

// BurnInfoCmd 获取下载视频命令数组
func BurnInfoCmd(burnInfoId, taskId uint, startTime, endTime *time.Time, tasks []*Task, saveFileTmpPath string) []*burninfocmd.BurnInfoCmd {
	// 查询下载文件名称
	videoName := burnsetting.GetBurnVideoName(taskId)

	var burnInfoCmds []*burninfocmd.BurnInfoCmd

	for index, task := range tasks {
		// 每个任务生成文件名称下标
		videoName = videoName + "-" + strconv.Itoa(index)
		// 比较 参数的开始时间 和 任务的开始时间 大小
		ss, duration := helper.CalculatingTime(startTime, endTime, task.FfmpegSaveStartTime, task.FfmpegSaveCloseTime)
		var mp4SavePath string
		if strings.LastIndex(saveFileTmpPath, "\\") == len(saveFileTmpPath)-1 {
			mp4SavePath = saveFileTmpPath + videoName + consts.SplitFileName
		} else {
			mp4SavePath = path.Join(saveFileTmpPath, videoName+consts.SplitFileName)
		}

		// 下载视频的ffmpeg命令构建
		cmdArgs, cmd := helper.GetSaveFileCmd(ss, task.M3u8Url, duration, mp4SavePath)

		// 存储子任务
		burnInfoCmd := burninfocmd.BurnInfoCmd{
			FfmpegCmd:     cmd,
			FfmpegCmdArgs: cmdArgs,
			BurnInfoID:    burnInfoId,
		}
		mysql.Instance.Create(&burnInfoCmd)

		burnInfoCmds = append(burnInfoCmds, &burnInfoCmd)
	}

	return burnInfoCmds
}

// GetDownloadCmds 获取下载视频命令数组
func GetDownloadCmds(videoName string, startTime, endTime *time.Time, tasks []*Task, fileSavePath string) []string {
	var cmds []string
	for index, task := range tasks {
		// 每个任务生成文件名称下标
		videoName = videoName + "-" + strconv.Itoa(index)
		// 比较 参数的开始时间 和 任务的开始时间 大小
		ss, duration := helper.CalculatingTime(startTime, endTime, task.FfmpegSaveStartTime, task.FfmpegSaveCloseTime)
		var mp4SavePath string
		if strings.LastIndex(fileSavePath, "\\") == len(fileSavePath)-1 {
			mp4SavePath = fileSavePath + videoName + consts.SplitFileName
		} else {
			mp4SavePath = path.Join(fileSavePath, videoName+consts.SplitFileName)
		}

		// 下载视频的ffmpeg命令构建
		_, cmd := helper.GetSaveFileCmd(ss, task.M3u8Url, duration, mp4SavePath)

		cmds = append(cmds, cmd)
	}
	return cmds
}
