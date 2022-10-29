package api

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/burninfo"
	"github.com/PittYao/stream_burn/internal/model/roommix3"
	"github.com/PittYao/stream_burn/internal/model/stream"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BurnMix3Video godoc
// @Summary 三合一
// @Tags 下载
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param mix3Req body dto.BurnMix3VideoDTO true " "
// @Router /api/v1/burnMixVideo [post]
func BurnMix3Video(c *gin.Context) {
	var mix3Req dto.BurnMix3VideoDTO
	err := c.ShouldBindJSON(&mix3Req)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验参数
	startTime, endTime, savFileTmpPath, done := CheckReq(c, mix3Req.StartTime, mix3Req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	mix3s := roommix3.QueryMix3File(mix3Req)
	if len(mix3s) == 0 {
		log.L.Error("没有查询到存在视频任务", zap.Any("data", mix3Req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 存储下载根任务
	burnInfo := burninfo.BurnInfo{
		BurnSettingID:   mix3Req.TaskId,
		BurnType:        consts.Mix3,
		CallbackUrl:     mix3Req.CallbackUrl,
		UndoneNum:       len(mix3s),
		StartTime:       startTime,
		EndTime:         endTime,
		SaveFileTmpPath: savFileTmpPath,
		OdaSavePath:     mix3Req.OdaSavePath,
		Uuid:            helper.RandomStr(),
	}
	mysql.Instance.Create(&burnInfo)
	// 获取下载视频命令行
	tasks := roommix3.ModelToTasks(mix3s)
	burnInfoCmds := stream.BurnInfoCmd(burnInfo.ID, mix3Req.TaskId, startTime, endTime, tasks, savFileTmpPath)
	// 开始下载所有任务
	for _, cmd := range burnInfoCmds {
		go cmd.DownloadVideo()
	}
	response.OKMsg(c, "开始下载视频", map[string]interface{}{
		"taskId": burnInfo.Uuid,
	})
}
