package publicrecordone

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

type PublicRecordOne struct {
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
	M3u8Url                  string     `json:"m3u8Url"` //m3u8地址
	FileRecentTime           *time.Time // 最新生成文件的时间
	TsFile                   string     // 最新ts文件地址
	DisuseAt                 *time.Time // 淘汰的时间，过期的文件可以被删除
}

// --- orm --- //

// Add 插入单个流任务
func (r *PublicRecordOne) Add() error {
	create := mysql.Instance.Create(r)
	if create.Error != nil {
		log.L.Error("新增公区转流任务失败", zap.Error(create.Error))
		return errors.New(" 新增公区转流任务失败")
	}

	return nil
}

// Update 更新
func (r *PublicRecordOne) Update() error {
	save := mysql.Instance.Save(&r)
	if save.Error != nil {
		log.L.Error("更新失败", zap.Error(save.Error))
		return errors.New("更新失败")
	}
	return save.Error

}

// GetById id查询
func GetById(id uint) (*PublicRecordOne, error) {
	var publicRecordOne PublicRecordOne
	mysql.Instance.First(&publicRecordOne, id)

	if publicRecordOne.ID == 0 {
		log.L.Error("DB中没有查询到该公区单画面任务", zap.Uint("id", id))
		return nil, errors.New("DB中没有查询到该公区单画面任务")
	}

	return &publicRecordOne, nil
}

// Delete 删除
func (r *PublicRecordOne) Delete() error {
	save := mysql.Instance.Delete(&r)
	if save.Error != nil {
		log.L.Error("公区 删除失败", zap.Error(save.Error))
		return errors.New("公区 删除失败")
	}
	return save.Error

}

// QueryPublicOneFile  查询房间单画面视频任务
func QueryPublicOneFile(burnSingleVideoDTO dto.BurnSingleVideoDTO) []*PublicRecordOne {
	var ones []*PublicRecordOne

	var middle []*PublicRecordOne
	var include []*PublicRecordOne
	var left []*PublicRecordOne
	var right []*PublicRecordOne
	var ing []*PublicRecordOne

	rtspUrl := helper.EncodeRtspUrl(burnSingleVideoDTO.RtspUrl)
	startTime := burnSingleVideoDTO.StartTime
	endTime := burnSingleVideoDTO.EndTime

	// 查询已经结束的任务 和 异常结束的任务能否满足查询条件
	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0 and ts_file != 0",
		rtspUrl, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&middle)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time > ? and ffmpeg_save_close_time < ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0 and ts_file != 0",
		rtspUrl, startTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&include)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time > ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time >= ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0 and ts_file != 0",
		rtspUrl, startTime, endTime, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&left)

	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ? and ffmpeg_save_close_time < ? and ffmpeg_save_close_time > ? and ffmpeg_save_state != ? and m3u8_url is not null AND LENGTH(trim(m3u8_url))>0 and ts_file != 0",
		rtspUrl, startTime, endTime, startTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&right)

	// 查询是否有正在进行的任务能满足查询条件
	mysql.Instance.Where("rtsp_url = ? and ffmpeg_save_start_time <= ?  and ffmpeg_save_state = ?",
		rtspUrl, endTime, consts.RunIng).Order("ffmpeg_save_start_time asc").Find(&ing)

	ones = append(ones, middle...)
	ones = append(ones, include...)
	ones = append(ones, left...)
	ones = append(ones, right...)

	if len(ones) != 0 {
		// 处理已经结束的任务
		for i := 0; i < len(ones); i++ {
			one := ones[i]
			if one.TsFile == "" {
				log.L.Sugar().Error("任务tsFile为空,任务id是:%d", one.ID)
				continue
			}
			// 校验m3u8地址是否可用
			err := httpclient.CheckM3u8Available(one.M3u8Url)
			if err != nil {
				log.L.Sugar().Error("任务m3u8Url不可用,m3u8Url:%s", one.M3u8Url)
				ones = append(ones[:i], ones[i+1:]...)
				i--
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
			copyUrl := helper.RedirectUrlBuilder(one.Ip, consts.PublicPort, fmt.Sprintf("/%s%s", consts.PublicSingle, consts.CopyM3u8))
			err, resp := httpclient.CopyM3u8HttpClient(copyUrl, one.ID)
			if err != nil {
				continue
			}

			m3u8TempUrl := fmt.Sprintf("%v", resp.Data)
			log.L.Info("获取临时m3u8文件成功", zap.String("url", copyUrl), zap.Any("m3u8TempUrl", m3u8TempUrl))
			one.M3u8Url = m3u8TempUrl

			now := time.Now()
			one.FfmpegSaveCloseTime = &now

			ones = append(ones, one)

		}
	}

	return ones

}

func (r *PublicRecordOne) ModelToTask() *stream.Task {
	task := stream.Task{}
	copier.Copy(&task, r)
	return &task
}

func ModelToTasks(ones []*PublicRecordOne) []*stream.Task {
	var tasks []*stream.Task

	for _, one := range ones {
		task := one.ModelToTask()
		tasks = append(tasks, task)
	}

	return tasks
}
