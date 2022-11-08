package roommix4

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

type RoomMix4 struct {
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
func (r *RoomMix4) Add() error {
	create := mysql.Instance.Create(r)
	if create.Error != nil {
		log.L.Error("RoomMix4 新增转流任务失败", zap.Error(create.Error))
		return errors.New(" 新增转流任务失败")
	}
	return nil
}

// Update 更新
func (r *RoomMix4) Update() error {
	save := mysql.Instance.Save(&r)
	if save.Error != nil {
		log.L.Error("RoomMix4 更新失败", zap.Error(save.Error))
		return errors.New("RoomMix4 更新失败")
	}
	return save.Error

}

// Delete
func (r *RoomMix4) Delete() error {
	save := mysql.Instance.Delete(&r)
	if save.Error != nil {
		log.L.Error("RoomMix4 删除失败", zap.Error(save.Error))
		return errors.New("RoomMix4 删除失败")
	}
	return save.Error

}

// GetById id查询
func GetById(id uint) (*RoomMix4, error) {
	var RoomMix4 RoomMix4
	mysql.Instance.First(&RoomMix4, id)

	if RoomMix4.ID == 0 {
		log.L.Sugar().Errorf("没有查询到该四合一画面任务 id:%d", id)
		return nil, errors.New("没有查询到该四合一画面任务")
	}

	return &RoomMix4, nil
}

// QueryMix4File 根据时间查询时间内的任务
func QueryMix4File(burnMixVideoDTO dto.BurnMixVideo4To1DTO) []*RoomMix4 {
	var mix4s []*RoomMix4

	var middle []*RoomMix4
	var include []*RoomMix4
	var left []*RoomMix4
	var right []*RoomMix4
	var mixIng []*RoomMix4

	rtspUrlMiddle := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlMiddle)
	rtspUrlSmallOne := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlSmallOne)
	rtspUrlSmallTwo := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlSmallTwo)
	rtspUrlSmallThree := helper.EncodeRtspUrl(burnMixVideoDTO.RtspUrlSmallThree)
	temperature := burnMixVideoDTO.Temperature
	startTime := burnMixVideoDTO.StartTime
	endTime := burnMixVideoDTO.EndTime

	// 查询已经结束的任务 和 异常结束的任务能否满足查询条件
	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_small_one = ? and rtsp_url_small_two = ? and rtsp_url_small_three = ? and temperature = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlSmallOne, rtspUrlSmallTwo, rtspUrlSmallThree, temperature, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&middle)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_small_one = ? and rtsp_url_small_two = ? and rtsp_url_small_three = ?  and temperature = ?  and ffmpeg_save_start_time > ? and ffmpeg_save_close_time < ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlSmallOne, rtspUrlSmallTwo, rtspUrlSmallThree, temperature, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&include)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_small_one = ? and rtsp_url_small_two = ? and rtsp_url_small_three = ?  and temperature = ?  and ffmpeg_save_start_time > ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlSmallOne, rtspUrlSmallTwo, rtspUrlSmallThree, temperature, startTime, endTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&left)

	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_small_one = ? and rtsp_url_small_two = ? and rtsp_url_small_three = ? and temperature = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time < ? and ffmpeg_save_close_time > ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0",
		rtspUrlMiddle, rtspUrlSmallOne, rtspUrlSmallTwo, rtspUrlSmallThree, temperature, startTime, endTime, startTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&right)

	// 查询是否有正在进行的任务能满足查询条件
	mysql.Instance.Where("rtsp_url_middle = ? and rtsp_url_small_one = ? and rtsp_url_small_two = ? and rtsp_url_small_three = ? and temperature = ? and ffmpeg_save_start_time <= ?  and ffmpeg_save_state = ?",
		rtspUrlMiddle, rtspUrlSmallOne, rtspUrlSmallTwo, rtspUrlSmallThree, temperature, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&mixIng)

	mix4s = append(mix4s, middle...)
	mix4s = append(mix4s, include...)
	mix4s = append(mix4s, left...)
	mix4s = append(mix4s, right...)

	if len(mix4s) != 0 {
		// 处理已经结束的任务
		for i := 0; i < len(mix4s); i++ {
			mix4 := mix4s[i]
			if mix4.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", mix4.ID)
				continue
			}
			// 校验m3u8地址是否可用
			err := httpclient.CheckM3u8Available(mix4.M3u8Url)
			if err != nil {
				log.L.Sugar().Error("任务m3u8Url不可用,m3u8Url:%s", mix4.M3u8Url)
				mix4s = append(mix4s[:i], mix4s[i+1:]...)
				i--
			}

		}
	}

	if len(mixIng) != 0 {
		// 处理正在运行的任务 获取临时m3u8文件
		for _, mix4 := range mixIng {
			if mix4.Ip == "" {
				log.L.Info("该任务没有服务器ip", zap.Any("mix4", mix4))
				continue
			}

			if mix4.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", mix4.ID)
				continue
			}

			// 请求远端获取临时文件
			copyUrl := helper.RedirectUrlBuilder(mix4.Ip, consts.Mix4Port, fmt.Sprintf("/%s%s", consts.Mix4, consts.CopyM3u8))
			err, resp := httpclient.CopyM3u8HttpClient(copyUrl, mix4.ID)
			if err != nil {
				continue
			}

			m3u8TempUrl := fmt.Sprintf("%v", resp.Data)
			log.L.Info("获取临时m3u8文件成功", zap.String("url", copyUrl), zap.Any("m3u8TempUrl", m3u8TempUrl))
			mix4.M3u8Url = m3u8TempUrl

			now := time.Now()
			mix4.FfmpegSaveCloseTime = &now

			mix4s = append(mix4s, mix4)

		}
	}

	return mix4s
}

func (r *RoomMix4) ModelToTask() *stream.Task {
	task := stream.Task{}
	copier.Copy(&task, r)
	return &task
}

func ModelToTasks(mix4s []*RoomMix4) []*stream.Task {
	var tasks []*stream.Task

	for _, mix4 := range mix4s {
		task := mix4.ModelToTask()
		tasks = append(tasks, task)
	}

	return tasks
}
