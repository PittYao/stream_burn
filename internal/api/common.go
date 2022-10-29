package api

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/helper"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"time"
)

// ValidatedReqTime 校验时间参数
func ValidatedReqTime(c *gin.Context, startTimeStr, endTimeStr string) (*time.Time, *time.Time, bool) {
	err, b, startTime, endTime := helper.TimeCompare(startTimeStr, endTimeStr)
	if err != nil {
		log.L.Info("开始时间或结束时间格式不正确 ", zap.String("startTime", startTimeStr), zap.String("endTime", endTimeStr))
		response.Err(c, "开始时间或结束时间格式不正确")
		return nil, nil, true
	}

	if !b {
		log.L.Info("开始时间不能比结束时间大", zap.String("startTime", startTimeStr), zap.String("endTime", endTimeStr))
		response.Err(c, "开始时间不能比结束时间大")
		return nil, nil, true
	}
	return startTime, endTime, false
}

// MkdirForTemp 创建下载文件临时文件夹
func MkdirForTemp(c *gin.Context) (string, bool) {
	saveFileTempPath := config.C.Burn.TmpPath + helper.RandomStr()
	err := os.MkdirAll(saveFileTempPath, os.ModePerm)
	if err != nil {
		log.L.Info("创建存储文件夹异常", zap.String("filePath", saveFileTempPath))
		response.Err(c, "创建存储文件夹异常")
		return "", true
	}
	return saveFileTempPath, false
}

func CheckReq(c *gin.Context, startTimeStr, endTimeStr string) (*time.Time, *time.Time, string, bool) {
	// 校验时间参数合法性
	startTime, endTime, done := ValidatedReqTime(c, startTimeStr, endTimeStr)
	if done {
		return nil, nil, "", true
	}
	// 创建临时文件夹 用于存放下载的文件
	savFileTmpPath, done := MkdirForTemp(c)
	if done {
		return nil, nil, "", true
	}

	return startTime, endTime, savFileTmpPath, false
}
