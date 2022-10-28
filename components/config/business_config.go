package config

import "time"

type Video struct {
	Type          string        `yaml:"type"`
	SaveDisk      string        `yaml:"saveDisk"`
	SaveDir       string        `yaml:"saveDir"`
	SaveDay       int           `yaml:"saveDay"`
	M3u8MaxTime   time.Duration `yaml:"m3u8MaxTime"`
	DiskConnect   bool          `yaml:"-"`
	TotalCapacity string        `yaml:"-"`
	UsedCapacity  string        `yaml:"-"`
	FreeCapacity  string        `yaml:"-"`
	UsedPercent   string        `yaml:"-"`
	DiskCapacity  int64         `yaml:"diskCapacity"` // 磁盘剩余容量达到diskCapacity时就删除历史文件 单位G
	DeleteDay     int64         `yaml:"deleteDay"`
}

func checkConfigAttribute() {
	if C.Video.Type == "" {
		panic("配置文件中video.type为空,不能启动服务")
	}

	if C.Video.SaveDisk == "" {
		panic("配置文件中video.saveDisk为空,不能启动服务")
	}

	if C.Video.SaveDir == "" {
		panic("配置文件中video.saveDir为空,不能启动服务")
	}
}
