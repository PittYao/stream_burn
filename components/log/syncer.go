package log

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"os"
)

// WriteSyncer 利用lumberjack库做日志分割
func WriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    config.C.Log.MaxSize,    // 单文件最大
		MaxBackups: config.C.Log.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     config.C.Log.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.C.Log.Compress,   // 是否压缩
		LocalTime:  config.C.Log.Localtime,  // 压缩文件名是否使用 localtime
	}
	// 只在控制台输出日志
	if config.C.Log.Format == consts.LogFormatConsole {
		return zapcore.AddSync(os.Stdout)
	}
	// 日志输出到控制台和文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
