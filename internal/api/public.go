package api

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/burninfo"
	"github.com/PittYao/stream_burn/internal/model/publicrecordone"
	"github.com/PittYao/stream_burn/internal/model/stream"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BurnPublicVideo godoc
// @Summary 下载公区
// @Tags 公区
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param publicReq body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/burnOtherSingleVideo [post]
func BurnPublicVideo(c *gin.Context) {
	var publicReq dto.BurnSingleVideoDTO
	if err := c.ShouldBindJSON(&publicReq); err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验参数
	startTime, endTime, savFileTmpPath, done := CheckReq(c, publicReq.StartTime, publicReq.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	oneFiles := publicrecordone.QueryPublicOneFile(publicReq)
	if len(oneFiles) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("data", publicReq))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 存储下载根任务
	burnInfo := burninfo.BurnInfo{
		BurnSettingID:   publicReq.TaskId,
		BurnType:        consts.PublicSingle,
		CallbackUrl:     publicReq.CallbackUrl,
		UndoneNum:       len(oneFiles),
		StartTime:       startTime,
		EndTime:         endTime,
		SaveFileTmpPath: savFileTmpPath,
		OdaSavePath:     publicReq.OdaSavePath,
		Uuid:            helper.RandomStr(),
	}
	mysql.Instance.Create(&burnInfo)
	// 获取下载视频命令行
	tasks := publicrecordone.ModelToTasks(oneFiles)
	burnInfoCmds := stream.BurnInfoCmd(burnInfo.ID, publicReq.TaskId, startTime, endTime, tasks, savFileTmpPath)
	// 开始下载所有任务
	for _, cmd := range burnInfoCmds {
		go cmd.DownloadVideo()
	}
	response.OKMsg(c, "开始下载视频", map[string]interface{}{
		"taskId": burnInfo.Uuid,
	})
}
