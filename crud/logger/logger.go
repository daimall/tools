package logger

import (
	"os"
	"path/filepath"

	"github.com/daimall/tools/crud/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func init() {
	if Logger != nil {
		return
	}
	logPath := common.GetPath([]string{"logs"})
	// 设置文件写入器的级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})

	// 创建日志核心
	core := zapcore.NewTee(
		newCore(filepath.Join(logPath, "app.log"), zapcore.DebugLevel),
		newCore(filepath.Join(logPath, "app.error.log"), zapcore.ErrorLevel),
		newCore(filepath.Join(logPath, "app.info.log"), infoLevel),
		newCore(filepath.Join(logPath, "app.warn.log"), warnLevel),
		newCore(filepath.Join(logPath, "app.debug.log"), debugLevel),
	)

	// 创建日志记录器
	Logger = zap.New(core, zap.ErrorOutput(zapcore.AddSync(os.Stderr)))
}

// NewCore creates a Core that writes logs to a WriteSyncer.
func newCore(logPath string, enab zapcore.LevelEnabler) zapcore.Core {
	// 创建 Lumberjack 配置
	viper.SetDefault("log.maxFileSize", 20) // 日志文件的最大大小（MB）
	viper.SetDefault("log.maxBackups", 5)   // 保留的旧日志文件的最大数量
	viper.SetDefault("log.maxAge", 5)       // 旧日志文件的最大保留天数
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    viper.GetInt("log.maxFileSize"),
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     viper.GetInt("log.maxAge"),
	}
	// 创建日志编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(lumberjackLogger), enab)
}
