package roommix3

import (
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/httpclient"
	"github.com/PittYao/stream_burn/internal/model/stream"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type RoomMix3 struct {
	gorm.Model
	RtspUrlMiddle            string
	RtspUrlLeft              string
	RtspUrlRight             string
	Temperature              string
	RoomName                 string
	Ip                       string
	Port                     string
	SavePath                 string
	FileRecentTime           *time.Time
	FfmpegTransformState     int
	FfmpegTransformCmd       string
	FfmpegTransformErrorMsg  string
	FfmpegTransformStartTime *time.Time
	FfmpegTransformCloseTime *time.Time
	FfmpegSaveState          int
	FfmpegSaveCmd            string
	FfmpegSaveErrorMsg       string
	FfmpegSaveStartTime      *time.Time
	FfmpegSaveCloseTime      *time.Time
	FfmpegStateLog           string
	TsFile                   string
	RebootRootId             uint
	RebootParentId           uint
	DisuseAt                 *time.Time
	M3u8Url                  string `json:"m3u8Url"`
}

// --- orm --- //

// Add 插入任务
func (r *RoomMix3) Add() error {
	create := mysql.Instance.Create(r)
	if create.Error != nil {
		log.L.Error("RoomMix3 新增mix3转流任务失败", zap.Error(create.Error))
		return errors.New(" 新增mix3转流任务失败")
	}
	//新增mix3转流任务失败   {"error": "Error 1292: Incorrect datetime value: '0000-00-00' for column 'file_recent_time' at row 1"}
	return nil
}

// Update 更新
func (r *RoomMix3) Update() error {
	save := mysql.Instance.Save(&r)
	if save.Error != nil {
		log.L.Error("RoomMix3 更新失败", zap.Error(save.Error))
		return errors.New("RoomMix3 更新失败")
	}
	return save.Error

}

// Delete
func (r *RoomMix3) Delete() error {
	save := mysql.Instance.Delete(&r)
	if save.Error != nil {
		log.L.Error("RoomMix3 删除失败", zap.Error(save.Error))
		return errors.New("RoomMix3 删除失败")
	}
	return save.Error

}

// GetById id查询
func GetById(id uint) (*RoomMix3, error) {
	var roomMix3 RoomMix3
	mysql.Instance.First(&roomMix3, id)

	if roomMix3.ID == 0 {
		log.L.Sugar().Errorf("没有查询到该3合一画面任务 id:%d", id)
		return nil, errors.New("没有查询到该3合一画面任务")
	}

	return &roomMix3, nil
}

// QueryMix3File 根据时间查询时间内的任务
func QueryMix3File(burnMixVideoDTO dto.BurnMix3VideoDTO) []*RoomMix3 {
	var mix3s []*RoomMix3

	var middle []*RoomMix3
	var include []*RoomMix3
	var left []*RoomMix3
	var right []*RoomMix3
	var mixIng []*RoomMix3

	rtspUrlMiddle := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlMiddle)
	rtspUrlLeft := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlLeft)
	rtspUrlRight := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlRight)
	temperature := burnMixVideoDTO.Temperature
	startTime := burnMixVideoDTO.StartTime
	endTime := burnMixVideoDTO.EndTime

	// 查询已经结束的任务 和 异常结束的任务能否满足查询条件
	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_left = ? and rtsp_url_right = ? and temperature = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlLeft, rtspUrlRight, temperature, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&middle)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_left = ? and rtsp_url_right = ? and temperature = ?  and ffmpeg_save_start_time > ? and ffmpeg_save_close_time < ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlLeft, rtspUrlRight, temperature, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&include)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_left = ? and rtsp_url_right = ? and temperature = ?  and ffmpeg_save_start_time > ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlLeft, rtspUrlRight, temperature, startTime, endTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&left)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_left = ? and rtsp_url_right = ? and temperature = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time < ? and ffmpeg_save_close_time > ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlLeft, rtspUrlRight, temperature, startTime, endTime, startTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&right)

	// 查询是否有正在进行的任务能满足查询条件
	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_left = ? and rtsp_url_right = ? and temperature = ? and ffmpeg_save_start_time <= ?  and ffmpeg_save_state = ?",
		rtspUrlMiddle, rtspUrlLeft, rtspUrlRight, temperature, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&mixIng)

	mix3s = append(mix3s, middle...)
	mix3s = append(mix3s, include...)
	mix3s = append(mix3s, left...)
	mix3s = append(mix3s, right...)

	if len(mix3s) != 0 {
		// 处理已经结束的任务
		for i := 0; i < len(mix3s); i++ {
			mix3 := mix3s[i]
			if mix3.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", mix3.ID)
				continue
			}
			// 校验m3u8地址是否可用
			err := httpclient.CheckM3u8Available(mix3.M3u8Url)
			if err != nil {
				log.L.Sugar().Error("任务m3u8Url不可用,m3u8Url:%s", mix3.M3u8Url)
				mix3s = append(mix3s[:i], mix3s[i+1:]...)
				i--
			}

		}
	}

	if len(mixIng) != 0 {
		// 处理正在运行的任务 获取临时m3u8文件
		for _, mix3 := range mixIng {
			if mix3.Ip == "" {
				log.L.Info("该任务没有服务器ip", zap.Any("mix3", mix3))
				continue
			}

			if mix3.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", mix3.ID)
				continue
			}

			// 请求远端获取临时文件
			copyUrl := helper.RedirectUrlBuilder(mix3.Ip, consts.Mix3Port, fmt.Sprintf("/%s%s", consts.Mix3, consts.CopyM3u8))
			err, resp := httpclient.CopyM3u8HttpClient(copyUrl, mix3.ID)
			if err != nil {
				continue
			}

			m3u8TempUrl := fmt.Sprintf("%v", resp.Data)
			log.L.Info("获取临时m3u8文件成功", zap.String("url", copyUrl), zap.Any("m3u8TempUrl", m3u8TempUrl))
			mix3.M3u8Url = m3u8TempUrl

			now := time.Now()
			mix3.FfmpegSaveCloseTime = &now

			mix3s = append(mix3s, mix3)

		}
	}

	return mix3s
}

func (r *RoomMix3) ModelToTask() *stream.Task {
	task := stream.Task{}
	copier.Copy(&task, r)
	return &task
}

func ModelToTasks(mix3s []*RoomMix3) []*stream.Task {
	var tasks []*stream.Task

	for _, mix3 := range mix3s {
		task := mix3.ModelToTask()
		tasks = append(tasks, task)
	}

	return tasks
}
