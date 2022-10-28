package components

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/components/swagger"
)

func Init() {
	config.Load()
	log.Load()
	swagger.Load()
	mysql.Load()
	log.L.Info("项目初始化配置完成")
}
