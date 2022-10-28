package burninfo

import (
	"errors"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type BurnInfo struct {
	gorm.Model
	BurnSettingID   uint
	BurnType        string
	CallbackUrl     string
	CallbackStatus  int64
	UndoneNum       int
	StartTime       *time.Time
	EndTime         *time.Time
	SaveFileTmpPath string
	OdaSavePath     string
	Uuid            string
}

func GetById(id uint) BurnInfo {
	var burnInfo BurnInfo
	mysql.Instance.Where("id = ?", id).First(&burnInfo)
	return burnInfo
}

func (b *BurnInfo) Update() error {
	save := mysql.Instance.Save(&b)
	if save.Error != nil {
		log.L.Error("BurnInfo 更新失败", zap.Error(save.Error))
		return errors.New("BurnInfo 更新失败")
	}
	return save.Error

}

// ReduceUndoneNum 减少未完成的任务数量
func (b *BurnInfo) ReduceUndoneNum(stepSize int) {
	b.UndoneNum = b.UndoneNum - stepSize
	mysql.Instance.Save(b)
}
