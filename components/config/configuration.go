package config

import (
	"fmt"
	"github.com/duke-git/lancet/fileutil"
	"github.com/spf13/viper"
	"os"
)

const (
	Profile        = "server.Profile"
	DevProfile     = "dev"
	configBaseFile = "app.yml"
	configFileName = "app"
	configType     = "yml"
	configBasePath = "/etc/"
)

var (
	baseConfigPath    string
	configDefaultPath string
)

func init() {
	projectPath, err := os.Getwd()
	if err != nil {
		panic("启动服务失败:" + err.Error())
	}
	configDefaultPath = projectPath + configBasePath
	baseConfigPath = configDefaultPath + configBaseFile
}

func Load() error {
	checkConfigFileExist(baseConfigPath)
	loadApplicationYml()
	loadProfileYml()
	checkConfigAttribute()
	return nil
}

// checkConfigFile 检测基础配置文件是否存在
func checkConfigFileExist(filepath string) {
	if !fileutil.IsExist(filepath) {
		panic("请检查配置文件" + filepath + "是否存在")
	}
}

// loadApplicationYml 读取基础配置文件
func loadApplicationYml() {
	v := viper.New()
	setViperAttribute(v, configFileName)
	// 将基础配置全部以默认配置写入
	configs := v.AllSettings()
	for k, v := range configs {
		viper.SetDefault(k, v)
	}
}

// loadProfileYml 获取指定配置文件
func loadProfileYml() {
	env := viper.GetString(Profile)
	if env == "" {
		panic(fmt.Sprintf("配置文件%s必须指定%s属性", configBaseFile, Profile))
	}

	// 校验指定配置文件是否存在
	configName := configFileName + "-" + env

	configFullName := getConfigFullName(configName)
	envConfigFilePath := configDefaultPath + configFullName
	checkConfigFileExist(envConfigFilePath)

	// 根据配置的env读取相应的配置信息
	setViperAttribute(viper.GetViper(), configName)

	C = &Config{}
	if err := viper.Unmarshal(C); err != nil {
		panic(fmt.Sprintf("解析配置文件异常err:%+v", err))
	}

}

// setViperAttribute 读取配置文件
func setViperAttribute(v *viper.Viper, configName string) {
	v.SetConfigName(configName)
	v.AddConfigPath(configDefaultPath)
	v.SetConfigType(configType)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("读取配置文件%s异常,err:%+v", getConfigFullName(configName), err.Error()))
	}
}

// getConfigFullName 获取配置文件全名
func getConfigFullName(configName string) string {
	return configName + "." + configType
}
