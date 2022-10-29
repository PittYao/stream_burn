package burnsetting

import (
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"gorm.io/gorm"
)

type BurnSetting struct {
	gorm.Model
	TaskNum        int    `gorm:"taskNum;comment:'刻录任务数'" json:"taskNum" binding:"required"`
	EncryptionType int    `gorm:"encryption_type;comment:'加密方式  1=不加密 2=密码加密 3=其他加密'"json:"encryptionType" binding:"required"`
	Password       string `gorm:"password;comment:'密码'" json:"password"`
	OdaSavePath    string `gorm:"oda_save_path;comment:'刻录文件存放路径'" json:"odaSavePath" binding:"required"`
	VideoName      string `gorm:"video_name;comment:'视频文件名称'" json:"videoName"`
	DoneTaskNum    int    `gorm:"done_task_num;comment:'已完成的任务数'" json:"doneTaskNum"`
}

// GetBurnVideoName 获取刻录视频的文件名
func GetBurnVideoName(taskId uint) string {
	var burnSetting BurnSetting
	mysql.Instance.First(&burnSetting, taskId)

	var videoName string
	if burnSetting.ID != 0 {
		videoName = burnSetting.VideoName
	}

	if videoName == "" {
		videoName = helper.RandomStr()
	}

	return videoName
}

// AddDoneTaskNum 完成数+1
func AddDoneTaskNum(id uint) (burnSetting BurnSetting) {
	mysql.Instance.Where("id = ?", id).First(&burnSetting)
	if burnSetting.ID == 0 {
		log.L.Info("查询BurnSetting不存在")
		return
	}
	burnSetting.DoneTaskNum += 1
	mysql.Instance.Save(&burnSetting)
	return
}
