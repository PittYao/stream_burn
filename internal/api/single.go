package api

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/burninfo"
	"github.com/PittYao/stream_burn/internal/model/roomrecordone"
	"github.com/PittYao/stream_burn/internal/model/stream"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BurnSingleVideo godoc
// @Summary 房间
// @Tags 下载
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param singleReq body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/burnSingleVideo [post]
func BurnSingleVideo(c *gin.Context) {
	var singleReq dto.BurnSingleVideoDTO
	err := c.ShouldBindJSON(&singleReq)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验参数
	startTime, endTime, savFileTmpPath, done := CheckReq(c, singleReq.StartTime, singleReq.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	oneFiles := roomrecordone.QuerySingleFile(singleReq)
	if len(oneFiles) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("data", singleReq))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 存储下载根任务
	burnInfo := burninfo.BurnInfo{
		BurnSettingID:   singleReq.TaskId,
		BurnType:        consts.Single,
		CallbackUrl:     singleReq.CallbackUrl,
		UndoneNum:       len(oneFiles),
		StartTime:       startTime,
		EndTime:         endTime,
		SaveFileTmpPath: savFileTmpPath,
		OdaSavePath:     singleReq.OdaSavePath,
		Uuid:            helper.RandomStr(),
	}
	mysql.Instance.Create(&burnInfo)
	// 获取下载视频命令行
	tasks := roomrecordone.ModelToTasks(oneFiles)
	burnInfoCmds := stream.BurnInfoCmd(burnInfo.ID, singleReq.TaskId, startTime, endTime, tasks, savFileTmpPath)
	// 开始下载所有任务
	for _, cmd := range burnInfoCmds {
		go cmd.DownloadVideo()
	}
	response.OKMsg(c, "开始下载视频", map[string]interface{}{
		"taskId": burnInfo.Uuid,
	})
}
