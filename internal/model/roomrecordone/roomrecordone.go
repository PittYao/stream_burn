package roomrecordone

import (
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/httpclient"
	"github.com/PittYao/stream_burn/internal/model/burninfocmd"
	"github.com/PittYao/stream_burn/internal/model/burnsetting"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RoomRecordOne struct {
	gorm.Model
	RtspUrl                  string
	Ip                       string
	Port                     string
	SavePath                 string
	FfmpegTransformState     int
	FfmpegTransformCmd       string
	FfmpegTransformErrorMsg  string
	FfmpegTransformStartTime *time.Time // 转流开始时间
	FfmpegTransformCloseTime *time.Time // 转流结束时间
	FfmpegSaveState          int
	FfmpegSaveCmd            string
	FfmpegSaveErrorMsg       string
	FfmpegSaveStartTime      *time.Time // 存流开始时间
	FfmpegSaveCloseTime      *time.Time // 存流结束时间
	FfmpegStateLog           string     // 流运行日志
	RebootRootId             uint       //重启任务的根id
	RebootParentId           uint       //重启任务的父id
	M3u8Url                  string     //m3u8地址
	FileRecentTime           *time.Time // 最新生成文件的时间
	TsFile                   string     // 最新ts文件
	DisuseAt                 *time.Time // 淘汰的时间，过期的文件可以被删除
}

// --- orm --- //

// Add 插入单个流任务
func (r *RoomRecordOne) Add() error {
	create := mysql.Instance.Create(r)
	if create.Error != nil {
		log.L.Error("新增转流任务失败", zap.Error(create.Error))
		return errors.New("新增转流任务失败")
	}

	return nil
}

// Update 更新
func (r *RoomRecordOne) Update() error {
	save := mysql.Instance.Save(&r)
	if save.Error != nil {
		log.L.Error("RoomRecordOne 更新失败", zap.Error(save.Error))
		return errors.New("RoomRecordOne 更新失败")
	}
	return save.Error

}

// Delete 删除
func (r *RoomRecordOne) Delete() error {
	save := mysql.Instance.Delete(&r)
	if save.Error != nil {
		log.L.Error("single 删除失败", zap.Error(save.Error))
		return errors.New("single 删除失败")
	}
	return save.Error

}

// GetById id查询
func GetById(id uint) (*RoomRecordOne, error) {
	var roomRecordOne RoomRecordOne
	mysql.Instance.First(&roomRecordOne, id)

	if roomRecordOne.ID == 0 {
		log.L.Error("DB中没有查询到该单画面任务", zap.Uint("id", id))
		return nil, errors.New("DB中没有查询到该单画面任务")
	}

	return &roomRecordOne, nil
}

// QuerySingleFile  查询房间单画面视频任务
func QuerySingleFile(burnSingleVideoDTO dto.BurnSingleVideoDTO) []*RoomRecordOne {
	var roomRecordOnes []*RoomRecordOne

	var middle []*RoomRecordOne
	var include []*RoomRecordOne
	var left []*RoomRecordOne
	var right []*RoomRecordOne
	var ing []*RoomRecordOne

	rtspUrl := helper.EncodeRtspUrl(burnSingleVideoDTO.RtspUrl)
	startTime := burnSingleVideoDTO.StartTime
	endTime := burnSingleVideoDTO.EndTime

	// 1.查询已经结束的任务 和 异常结束的任务能否满足查询条件
	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ?",
		rtspUrl, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&middle)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time > ? and ffmpeg_save_close_time < ? and ffmpeg_save_state != ?",
		rtspUrl, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&include)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time > ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ?",
		rtspUrl, startTime, endTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&left)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time < ? and ffmpeg_save_close_time > ? and ffmpeg_save_state != ?",
		rtspUrl, startTime, endTime, startTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&right)

	// 查询是否有正在进行的任务能满足查询条件
	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ?  and ffmpeg_save_state = ?",
		rtspUrl, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&ing)

	roomRecordOnes = append(roomRecordOnes, middle...)
	roomRecordOnes = append(roomRecordOnes, include...)
	roomRecordOnes = append(roomRecordOnes, left...)
	roomRecordOnes = append(roomRecordOnes, right...)

	if len(roomRecordOnes) != 0 {
		// 处理已经结束的任务
		for i, _ := range roomRecordOnes {
			one := roomRecordOnes[i]
			if one.M3u8Url == "" || one.TsFile == "" {
				log.L.Sugar().Error("任务m3u8Url或tsFile为空,任务id是:%d", one.ID)
				continue
			}
			// 校验m3u8地址是否可用
			err := httpclient.CheckM3u8Available(one.M3u8Url)
			if err != nil {
				log.L.Sugar().Error("任务m3u8Url不可用,m3u8Url:%s", one.M3u8Url)
				roomRecordOnes = append(roomRecordOnes[:i], roomRecordOnes[i+1:]...)
			}

		}
	}

	if len(ing) != 0 {
		// 处理正在运行的任务 获取临时m3u8文件
		for _, one := range ing {
			if one.Ip == "" {
				log.L.Info("该任务没有服务器ip", zap.Any("one", one))
				continue
			}

			if one.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", one.ID)
				continue
			}

			// 请求远端获取临时文件
			copyUrl := helper.RedirectUrlBuilder(one.Ip, consts.SinglePort, fmt.Sprintf("/%s%s", consts.Single, consts.CopyM3u8))
			err, resp := httpclient.CopyM3u8HttpClient(copyUrl, one.ID)
			if err != nil {
				continue
			}

			m3u8TempUrl := fmt.Sprintf("%v", resp.Data)
			log.L.Info("获取临时m3u8文件成功", zap.String("url", copyUrl), zap.Any("m3u8TempUrl", m3u8TempUrl))
			one.M3u8Url = m3u8TempUrl

			now := time.Now()
			one.FfmpegSaveCloseTime = &now

			roomRecordOnes = append(roomRecordOnes, one)

		}
	}

	return roomRecordOnes

}

// BuildRoomOneMp4 下载视频命令构建
func BuildRoomOneMp4(burnInfoId, taskId uint, startTime, endTime *time.Time, oneFiles []*RoomRecordOne, savFileTmpPath string) []*burninfocmd.BurnInfoCmd {
	//  查询下载文件名称
	videoName := burnsetting.GetBurnVideoName(taskId)

	var burnInfoCmds []*burninfocmd.BurnInfoCmd

	for index, oneFile := range oneFiles {
		// 每个任务生成文件名称下标
		videoName = videoName + "-" + strconv.Itoa(index)

		// 比较 参数的开始时间 和 任务的开始时间 大小
		ss, duration := helper.CalculatingTime(startTime, endTime, oneFile.FfmpegSaveStartTime, oneFile.FfmpegSaveCloseTime)

		mp4SavePath := savFileTmpPath + "/" + videoName + consts.SplitFileName

		// 下载视频的ffmpeg命令构建
		cmdArgs, cmd := helper.GetSaveFileCmd(ss, oneFile.M3u8Url, duration, mp4SavePath)

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
