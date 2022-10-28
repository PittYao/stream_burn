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
	"github.com/PittYao/stream_burn/internal/model/roomrecordone"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

// BurnSingleVideo godoc
// @Summary 下载房间单画面
// @Tags 房间单画面
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param burnSingleVideoDTO body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/burnSingleVideo [post]
func BurnSingleVideo(c *gin.Context) {
	var burnSingleVideoDTO dto.BurnSingleVideoDTO
	err := c.ShouldBindJSON(&burnSingleVideoDTO)
	if err != nil {
		response.Err(c, err.Error())
		return
	}

	// 校验时间参数
	err, b, startTime, endTime := helper.TimeCompare(burnSingleVideoDTO.StartTime, burnSingleVideoDTO.EndTime)
	if err != nil {
		log.L.Info("开始时间或结束时间格式不正确 ", zap.Any("data", burnSingleVideoDTO))
	}

	if !b {
		log.L.Info("开始时间不能比结束时间大", zap.Any("data", burnSingleVideoDTO))
		response.Err(c, "开始时间不能比结束时间大")
		return
	}

	// 查询m3u8文件列表
	oneFiles := roomrecordone.QuerySingleFile(burnSingleVideoDTO)
	if len(oneFiles) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("data", burnSingleVideoDTO))
		response.Err(c, "没有查询到存在视频任务")
		return
	}

	// 创建临时文件夹 用于存放文件
	savFileTmpPath := config.C.Burn.TmpPath + helper.RandomStr()
	err = os.MkdirAll(savFileTmpPath, os.ModePerm)
	if err != nil {
		log.L.Info("下载房间单画面视频时,创建存储文件夹异常", zap.String("filePath", savFileTmpPath))
		return
	}

	// 存储下载根任务
	burnInfo := burninfo.BurnInfo{
		BurnSettingID:   burnSingleVideoDTO.TaskId,
		BurnType:        consts.Single,
		CallbackUrl:     burnSingleVideoDTO.CallbackUrl,
		UndoneNum:       len(oneFiles),
		StartTime:       startTime,
		EndTime:         endTime,
		SaveFileTmpPath: savFileTmpPath,
		OdaSavePath:     burnSingleVideoDTO.OdaSavePath,
		Uuid:            helper.RandomStr(),
	}
	mysql.Instance.Create(&burnInfo)

	// 下载视频命令构建
	burnInfoCmds := roomrecordone.BuildRoomOneMp4(burnInfo.ID, burnSingleVideoDTO.TaskId, startTime, endTime, oneFiles, savFileTmpPath)

	// 开始下载所有任务
	for _, cmd := range burnInfoCmds {
		go cmd.FfmpegDownloadRoomOneVideo()
	}

	response.OKMsg(c, "开始下载视频成功", map[string]interface{}{
		"taskId": burnInfo.Uuid,
	})
}
