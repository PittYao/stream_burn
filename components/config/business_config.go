package config

import (
	"os"
)

type Burn struct {
	TmpPath      string `yaml:"tmpPath"`
	CopyFilePath string `yaml:"copyFilePath"`
}

func checkConfigAttribute() {
	if C.Burn.TmpPath == "" {
		panic("配置文件中burn.tmpPath为空,不能启动服务")
	}

	err := os.MkdirAll(C.Burn.TmpPath, os.ModePerm)
	if err != nil {
		panic("创建存储文件夹失败,检测配置文件中burn.tmpPath是否正确")
	}

	if C.Burn.CopyFilePath == "" {
		panic("配置文件中burn.copyFilePath为空,不能启动服务")
	}
}
