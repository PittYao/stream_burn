package burnfile

import (
	"errors"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BurnFile struct {
	gorm.Model     `json:"-"`
	FileUrl        string `json:"fileUrl" example:"https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fhbimg.b0.upaiyun.com%2Fc1a81a1ff39c10b5a93fb76cc6d6d857c163aafd43a48-J1sDec_fw658&refer=http%3A%2F%2Fhbimg.b0.upaiyun.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1654910136&t=29cff2e1390a9e6bf24ea635c9156b79"`
	CallbackUrl    string `json:"callbackUrl" example:"http://localhost:8010/api/v1/callback"`
	CallbackStatus int64  `json:"callbackStatus"`
	OdaSavePath    string `json:"odaSavePath" example:"D:/downloadVideo"`
	DoneStatus     int64  `json:"done_status"`
	BurnSettingID  uint   `json:"taskId" example:"1"`
	Uuid           string `json:"-"`
}

func (b *BurnFile) Update() error {
	save := mysql.Instance.Save(&b)
	if save.Error != nil {
		log.L.Error("BurnFile 更新失败", zap.Error(save.Error))
		return errors.New("BurnFile 更新失败")
	}
	return save.Error

}
