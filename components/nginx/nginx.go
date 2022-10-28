package nginx

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"os/exec"
	"path/filepath"
)

func Load() error {
	// 检测nginx端口是否已开启
	if err := helper.CheckPortRunning(consts.Localhost, consts.RtmpPort); err != nil {
		log.L.Sugar().Errorf("nginx未启动,开始启动")
		// 启动nginx
		err := RunNginx()
		if err != nil {
			log.L.Sugar().Errorf("nginx启动失败,err:%+v", err)
			panic("nginx启动失败")
		} else {
			log.L.Sugar().Info("nginx启动成功")
			return nil
		}
	}

	log.L.Info("nginx已启动")
	return nil
}

// RunNginx 启动nginx
func RunNginx() (err error) {
	dir, _ := filepath.Abs(filepath.Dir("./"))

	args := []string{
		"-c",
		"conf/nginx.conf",
	}

	name := dir + "/" + config.C.Nginx.LibPath + "/nginx"
	nginx := exec.Command(name, args...)
	nginx.Dir = dir + "/" + config.C.Nginx.LibPath

	return nginx.Start()
}
