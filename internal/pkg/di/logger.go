package di

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Level      string // debug, info, warn, error
	Path       string // 日志保存路径
	FileName   string // 日志文件名 (e.g., app.log)
	MaxSize    int    // 单个文件最大尺寸 (MB)
	MaxBackups int    // 保留旧文件最大个数
	MaxAge     int    // 保留旧文件最大天数
	Compress   bool   // 是否压缩
	Console    bool   // 是否输出到控制台
	LocalTime  bool   // 是否使用本地时间
}

// InitLogger 初始化 Logger
func InitLogger() *zap.Logger {
	cfg := &LogConfig{
		Level:      viper.GetString("log.level"),
		Path:       viper.GetString("log.path"),
		FileName:   viper.GetString("log.filename"),
		MaxSize:    viper.GetInt("log.maxSize"),
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     viper.GetInt("log.maxAge"),
		Compress:   viper.GetBool("log.compress"),
		Console:    viper.GetBool("log.console"),
		LocalTime:  viper.GetBool("log.localTime"),
	}

	if cfg.Path == "" {
		cfg.Path = "./logs"
	}
	if cfg.FileName == "" {
		cfg.FileName = "app.log"
	}

	// 创建日志目录
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		panic("无法创建日志目录: " + err.Error())
	}

	// 配置 Lumberjack
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, cfg.FileName),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  cfg.LocalTime,
	})

	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zap.InfoLevel // 默认级别
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	// 文件编码器配置
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 控制台编码器配置
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var cores []zapcore.Core

	// 添加文件输出 Core
	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderConfig),
		writeSyncer,
		atomicLevel,
	))

	// 添加控制台输出 Core
	if cfg.Console {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.Lock(os.Stdout),
			atomicLevel,
		))
	}

	// 组合 Core 并创建 Logger
	core := zapcore.NewTee(cores...)

	// AddCaller: 添加调用行号
	// AddStacktrace: 只有 Error 级别以上才打印堆栈
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 替换全局 Logger
	zap.ReplaceGlobals(logger)

	return logger
}
