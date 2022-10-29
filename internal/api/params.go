package api

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/publicrecordone"
	"github.com/PittYao/stream_burn/internal/model/roommix3"
	"github.com/PittYao/stream_burn/internal/model/roomrecordone"
	"github.com/PittYao/stream_burn/internal/model/stream"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BurnParams godoc
// @Summary 三合一
// @Tags 下载参数
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param burnMix3VideoDTO body dto.BurnMix3VideoDTO true " "
// @Router /api/v1/burnParams [post]
func BurnParams(c *gin.Context) {
	var mix3Req dto.BurnMix3VideoDTO
	err := c.ShouldBindJSON(&mix3Req)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	startTime, endTime, done := ValidatedReqTime(c, mix3Req.StartTime, mix3Req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	mix3s := roommix3.QueryMix3File(mix3Req)
	if len(mix3s) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", mix3Req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 下载视频命令构建
	tasks := roommix3.ModelToTasks(mix3s)
	cmds := stream.GetDownloadCmds(mix3Req.VideoName, startTime, endTime, tasks, mix3Req.FileSavePath)
	response.OKMsg(c, "获取ffmpeg参数成功", cmds)
}

// BurnSingleParams godoc
// @Summary 房间单画面
// @Tags 下载参数
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param singleReq body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/burnSingleParams [post]
func BurnSingleParams(c *gin.Context) {
	var singleReq dto.BurnSingleVideoDTO
	err := c.ShouldBindJSON(&singleReq)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	startTime, endTime, done := ValidatedReqTime(c, singleReq.StartTime, singleReq.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	ones := roomrecordone.QuerySingleFile(singleReq)
	if len(ones) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", singleReq))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 下载视频命令构建
	tasks := roomrecordone.ModelToTasks(ones)
	cmds := stream.GetDownloadCmds(singleReq.VideoName, startTime, endTime, tasks, singleReq.FileSavePath)
	response.OKMsg(c, "获取ffmpeg参数成功", cmds)
}

// BurnOtherSingleParams godoc
// @Summary 公区
// @Tags 下载参数
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param publicReq body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/burnOtherSingleParams [post]
func BurnOtherSingleParams(c *gin.Context) {
	var publicReq dto.BurnSingleVideoDTO
	err := c.ShouldBindJSON(&publicReq)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	startTime, endTime, done := ValidatedReqTime(c, publicReq.StartTime, publicReq.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	ones := publicrecordone.QueryPublicOneFile(publicReq)
	if len(ones) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", publicReq))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	// 下载视频命令构建
	tasks := publicrecordone.ModelToTasks(ones)
	cmds := stream.GetDownloadCmds(publicReq.VideoName, startTime, endTime, tasks, publicReq.FileSavePath)
	response.OKMsg(c, "获取ffmpeg参数成功", cmds)
}
