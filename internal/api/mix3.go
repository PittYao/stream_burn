package api

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/burninfo"
	"github.com/PittYao/stream_burn/internal/model/roommix3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

// BurnMix3Video godoc
// @Summary 下载三合一画面
// @Tags 三合一
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param burnMix3VideoDTO body dto.BurnMix3VideoDTO true " "
// @Router /api/v1/burnMixVideo [post]
func BurnMix3Video(c *gin.Context) {
	var burnMix3VideoDTO dto.BurnMix3VideoDTO
	err := c.ShouldBindJSON(&burnMix3VideoDTO)
	if err != nil {
		response.Err(c, err.Error())
		return
	}

	// 校验时间参数
	err, b, startTime, endTime := helper.TimeCompare(burnMix3VideoDTO.StartTime, burnMix3VideoDTO.EndTime)
	if err != nil {
		log.L.Error("开始时间或结束时间格式不正确 ", zap.Any("data", burnMix3VideoDTO))
	}

	if !b {
		log.L.Error("开始时间不能比结束时间大", zap.Any("data", burnMix3VideoDTO))
		response.Err(c, "开始时间不能比结束时间大")
		return
	}

	// 查询m3u8文件列表
	mix3s := roommix3.QueryMix3File(burnMix3VideoDTO)
	if len(mix3s) == 0 {
		log.L.Error("没有查询到存在视频任务", zap.Any("data", burnMix3VideoDTO))
		response.Err(c, "没有查询到存在视频任务")
		return
	}

	// 创建临时文件夹 用于存放文件
	savFileTmpPath := config.C.Burn.TmpPath + helper.RandomStr()
	err = os.MkdirAll(savFileTmpPath, os.ModePerm)
	if err != nil {
		log.L.Error("下载合成视频时,创建存储文件夹异常", zap.String("filePath", savFileTmpPath))
		return
	}

	// 存储下载根任务
	burnInfo := burninfo.BurnInfo{
		BurnSettingID:   burnMix3VideoDTO.TaskId,
		BurnType:        consts.Mix3,
		CallbackUrl:     burnMix3VideoDTO.CallbackUrl,
		UndoneNum:       len(mix3s),
		StartTime:       startTime,
		EndTime:         endTime,
		SaveFileTmpPath: savFileTmpPath,
		OdaSavePath:     burnMix3VideoDTO.OdaSavePath,
		Uuid:            helper.RandomStr(),
	}
	mysql.Instance.Create(&burnInfo)

	// 下载视频命令构建
	burnInfoCmds := roommix3.BuildMix3Mp4(burnInfo.ID, burnMix3VideoDTO.TaskId, startTime, endTime, mix3s, savFileTmpPath)

	// 开始下载所有任务
	for _, cmd := range burnInfoCmds {
		go cmd.FfmpegDownloadMixVideo()
	}

	response.OKMsg(c, "开始下载视频成功", map[string]interface{}{
		"taskId": burnInfo.Uuid,
	})
}
