package api

import (
	"errors"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/httpclient"
	"github.com/PittYao/stream_burn/internal/model/burnfile"
	"github.com/PittYao/stream_burn/internal/model/burnsetting"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"os"
)

// BurnFile godoc
// @Summary 文件
// @Tags 下载
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param fileReq body dto.BurnFileDTO true " "
// @Router /api/v1/burnFile [post]
func BurnFile(c *gin.Context) {
	var fileReq dto.BurnFileDTO
	err := c.ShouldBindJSON(&fileReq)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	// 存储下载文件任务信息
	var burnFile burnfile.BurnFile
	copier.Copy(&burnFile, fileReq)
	burnFile.Uuid = helper.RandomStr()
	mysql.Instance.Save(&burnFile)
	// 下载文件
	go func() {
		err = FileDownLoad(burnFile)
		// 更新刻录文件状态
		if err != nil {
			burnFile.DoneStatus = consts.RunIngError
		} else {
			burnFile.DoneStatus = consts.Success
		}

		// 回调业务端接口
		err = httpclient.CallBackFileHttpClient(burnFile.CallbackUrl, burnFile.Uuid, err)
		if err != nil {
			burnFile.CallbackStatus = consts.RunIngError
		} else {
			burnFile.CallbackStatus = consts.Success
		}

		burnFile.Update()
	}()
	// 响应
	response.OKMsg(c, "开启下载", map[string]interface{}{
		"taskId": burnFile.Uuid,
		"status": consts.Success,
	})

}

// FileDownLoad 下载文件
func FileDownLoad(burnFile burnfile.BurnFile) error {
	// 创建临时路径文件夹
	err := os.MkdirAll(config.C.Burn.TmpPath, os.ModePerm)
	if err != nil {
		log.L.Error("创建存储文件夹失败", zap.Error(err))
		return errors.New("创建存储文件夹失败")
	}
	// 访问url存储到临时文件夹下
	err, savePath := helper.DownloadFile2Path(burnFile.FileUrl, config.C.Burn.TmpPath)
	if err != nil {
		return err
	}
	// 剪切到oda刻盘存储路径
	err = os.MkdirAll(burnFile.OdaSavePath, os.ModePerm)
	err = helper.CopyDir(savePath, burnFile.OdaSavePath, true)
	if err != nil {
		log.L.Error("剪切文件失败", zap.String("savePath", savePath), zap.Error(err))
		return errors.New("剪切文件失败")
	}
	// 任务完成数+1
	burnsetting.AddDoneTaskNum(burnFile.BurnSettingID)
	return nil
}
