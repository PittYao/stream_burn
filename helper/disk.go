package helper

import (
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	human "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
)

// LoadDiskInfo 获取磁盘容量信息
func LoadDiskInfo(diskPath string) (totalCapacity, usedCapacity, freeCapacity, usedPercent string, err error) {
	usage, err := disk.Usage(diskPath)

	if usage == nil || err != nil {
		return "", "", "", "", errors.New("配置文件中video.saveDisk的硬盘不存在,请修改存在的盘符")
	}
	totalCapacity = human.Bytes(usage.Total)
	usedCapacity = human.Bytes(usage.Used)
	freeCapacity = human.Bytes(usage.Free)
	usedPercent = fmt.Sprintf("%2.f%%", usage.UsedPercent)

	return
}

// FreeDiskOverLimit 检测磁盘容量已达到预警容量
func FreeDiskOverLimit(diskPath string, limitCapacity int64) (bool, error) {
	usage, err := disk.Usage(diskPath)
	if err != nil {
		return false, err
	}

	limitCapacityStr := fmt.Sprintf("%d Gib", limitCapacity)
	bytes, err := human.ParseBytes(limitCapacityStr)
	if err != nil {
		return false, err
	}

	if usage.Free < bytes {
		free := human.Bytes(usage.Free)
		log.L.Warn("磁盘已达到预警容量", zap.String("当前剩余容量", free), zap.String("预警容量", limitCapacityStr))
		return true, err
	}

	return false, err
}

func CheckDiskOverLimit() (bool, error) {
	diskOverLimit, err := FreeDiskOverLimit(config.C.Video.SaveDisk, config.C.Video.DiskCapacity)
	if err != nil {
		log.L.Info("获取磁盘容量达到预警值", zap.String("disk", config.C.Video.SaveDisk))
		return false, errors.New("获取磁盘容量达到预警值")
	}

	return diskOverLimit, nil
}
