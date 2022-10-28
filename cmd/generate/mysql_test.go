package gen

import (
	"github.com/qmhball/db2gorm/gen"
	"testing"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/stream_burn?charset=utf8&parseTime=True&loc=Local"

// 确保 WritePath文件夹已存在

// 生成指定单表
func TestGenSingleTableStruct(t *testing.T) {
	tblName := "ffmpeg_templates"
	gen.GenerateOne(gen.GenConf{
		Dsn:       dsn,
		WritePath: "../../internal/model",
		Stdout:    false,
		Overwrite: false,
	}, tblName)
}

// 生成指定整个Db
func TestGenDBStruct(t *testing.T) {
	gen.GenerateAll(gen.GenConf{
		Dsn:       dsn,
		WritePath: "../../internal/model",
		Stdout:    false,
		Overwrite: true,
	})
}
