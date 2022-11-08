package api

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/roommix3"
	"github.com/PittYao/stream_burn/internal/model/roommix4"
	"github.com/PittYao/stream_burn/internal/model/roomrecordone"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListMix3Video godoc
// @Summary 三合一
// @Tags 回放
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param req body dto.BurnMix3VideoDTO true " "
// @Router /api/v1/monitor/web/mix [post]
func ListMix3Video(c *gin.Context) {
	var req dto.BurnMix3VideoDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	_, _, done := ValidatedReqTime(c, req.StartTime, req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	mix3s := roommix3.QueryMix3File(req)
	if len(mix3s) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	response.OKMsg(c, "查询成功", mix3s)
}

// ListMix4Video godoc
// @Summary 四合一
// @Tags 回放
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param req body dto.BurnMix3VideoDTO true " "
// @Router /api/v1/monitor/web/mix4 [post]
func ListMix4Video(c *gin.Context) {
	var req dto.BurnMixVideo4To1DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	_, _, done := ValidatedReqTime(c, req.StartTime, req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	mix4s := roommix4.QueryMix4File(req)
	if len(mix4s) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	response.OKMsg(c, "查询成功", mix4s)
}

// ListSingleVideo godoc
// @Summary 房间单画面
// @Tags 回放
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param req body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/monitor/web/single [post]
func ListSingleVideo(c *gin.Context) {
	var req dto.BurnSingleVideoDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	_, _, done := ValidatedReqTime(c, req.StartTime, req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	ones := roomrecordone.QuerySingleFile(req)
	if len(ones) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	response.OKMsg(c, "查询成功", ones)
}

// ListPublicSingleVideo godoc
// @Summary 公区
// @Tags 回放
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param req body dto.BurnSingleVideoDTO true " "
// @Router /api/v1/monitor/web/other [post]
func ListPublicSingleVideo(c *gin.Context) {
	var req dto.BurnSingleVideoDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err.Error())
		return
	}
	// 校验时间参数
	_, _, done := ValidatedReqTime(c, req.StartTime, req.EndTime)
	if done {
		return
	}
	// 查询m3u8文件列表
	ones := roomrecordone.QuerySingleFile(req)
	if len(ones) == 0 {
		log.L.Info("没有查询到存在视频任务", zap.Any("req", req))
		response.Err(c, "没有查询到存在视频任务")
		return
	}
	response.OKMsg(c, "查询成功", ones)
}
