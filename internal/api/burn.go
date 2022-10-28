package api

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/dto"
	"github.com/PittYao/stream_burn/internal/model/burnsetting"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"os"
)

// BurnTask godoc
// @Summary 创建刻录任务
// @Tags 刻录任务
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param burnSettingDTO body dto.BurnSettingDTO true " "
// @Router /api/v1/burnTask [post]
func BurnTask(c *gin.Context) {
	var burnSettingDTO dto.BurnSettingDTO
	err := c.ShouldBindJSON(&burnSettingDTO)
	if err != nil {
		response.Err(c, err.Error())
		return
	}

	copyOne := burnsetting.BurnSetting{}
	copier.Copy(&copyOne, &burnSettingDTO)

	// create dir
	err = os.MkdirAll(burnSettingDTO.OdaSavePath, os.ModePerm)
	if err != nil {
		response.Err(c, "文件夹创建失败,请检查文件夹路径")
	}

	// copy file to odaSavePath
	helper.CopyDir(config.C.Burn.CopyFilePath, burnSettingDTO.OdaSavePath, false)

	mysql.Instance.Save(&copyOne)

	response.OKMsg(c, "创建任务成功", map[string]interface{}{
		"taskId": copyOne.ID,
	})

}
