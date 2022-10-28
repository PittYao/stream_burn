package config

type Burn struct {
	TmpPath      string `yaml:"tmpPath"`
	CopyFilePath string `yaml:"copyFilePath"`
}

func checkConfigAttribute() {
	if C.Burn.TmpPath == "" {
		panic("配置文件中burn.tmpPath为空,不能启动服务")
	}

	if C.Burn.CopyFilePath == "" {
		panic("配置文件中burn.copyFilePath为空,不能启动服务")
	}
}
