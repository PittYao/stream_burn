package swagger

import (
	"bytes"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func Load() error {
	env := viper.GetString(config.Profile)
	// 开发环境生成swagger 不会马上生效需重启服务才会生效
	if env == config.DevProfile {
		dir, err := os.Getwd()
		if err != nil {
			log.L.Sugar().Errorf("生成swagger异常, err:%+v", err)
			panic("生成swagger异常")
		}

		swaggerBatPath := dir + "\\components\\swagger\\gen.bat"
		cmd := exec.Command("cmd.exe", "/C", swaggerBatPath)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			log.L.Sugar().Errorf("生成swagger异常, err:%s", out.String())
			panic("生成swagger异常")
		}

		log.L.Sugar().Info("生成swagger成功")
	}

	return nil
}
