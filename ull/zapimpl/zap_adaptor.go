package zapimpl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/llmuz/yggdrasill/ull/config"
)

const (
	// 配置日志的等级
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	levelFatal = "fatal"
)

// NewZapLogger 日志配置文件
// 如果需要配置 options, 请使用 logger.WithOptions() 方法进行配置
// zap.AddCaller()
// zap.AddCallerSkip(2)
// 添加额外的配置
func NewZapLogger(cfg *config.LogConfig) (logger *zap.Logger, err error) {
	var (
		defaultBackUp   = 200       // 保留日志的最大值
		defaultSize     = 1024      // 默认日志最大分割容量
		defaultAge      = 7         // 日志保留的最大天数
		defaultFileName = "biz.log" // 默认日志文件名
	)

	var handleErr = func(msg string) (logger *zap.Logger, err error) {
		return nil, errors.New(msg)
	}

	if cfg == nil {
		return handleErr("NewZapLogger couldn't be nil")
	}

	// info is default log level
	var logLevel zapcore.Level
	if logLevel, err = parseLogLevel(cfg.GetLevel()); err != nil {
		return handleErr(err.Error())
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		CallerKey:      "file",
		SkipLineEnding: false,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		StacktraceKey:  "stack",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.9999999"))
		}, // time format
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 10e6)
		}, // duration
	}

	// 将所有的日志文件输出到同一个文件
	bizLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	// 保留文件的最大数量
	var maxBackupSize = defaultBackUp
	if cfg.MaxBackup != nil {

		maxBackupSize = int(cfg.GetMaxBackup())
	}

	// 保留日志的最大天数
	var maxAge = defaultAge
	if cfg.MaxAge != nil {
		maxAge = int(cfg.GetMaxAge())
	}

	// 日志的最大值
	var maxSize = defaultSize
	if cfg.MaxSize != nil {
		maxSize = int(cfg.GetMaxSize())
	}

	if cfg.FileName == "" {
		cfg.FileName = defaultFileName
	}
	// writer
	bizWriter := getWriter(cfg.GetFileName(),
		maxBackupSize, maxAge, maxSize, cfg.GetCompress())

	// 判断输入日志的格式
	var zc zapcore.Core
	if cfg.GetJsonFormat() {
		zc = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(bizWriter), bizLevel)
	} else {
		zc = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(bizWriter), bizLevel)
	}
	// 配置多个输出文件
	core := zapcore.NewTee(
		zc,
	)

	// debug 日志级别是否输出到控制台
	if cfg.GetDebugModeOutputConsole() && (strings.ToLower(cfg.GetLevel()) == levelDebug) {
		//同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		core = zapcore.NewTee(core, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), logLevel))
	}

	return zap.New(core), nil
}

func getWriter(filename string, maxBackup, maxAge, maxSize int, compress bool) io.Writer {
	fmt.Printf("getWriter %s, maxBackup:%d, maxAge:%d, maxSize:%dmb, compress:%v\n", filename, maxBackup, maxAge, maxSize, compress)
	return &lumberjack.Logger{
		Filename:   filename,  // 文件名
		MaxSize:    maxSize,   // the file max size, unit is mb, if overflow the file will rotate
		MaxBackups: maxBackup, // 最大文件保留数, 超过就删除最老的日志文件
		MaxAge:     maxAge,    // 保留文件的最大天数
		Compress:   compress,  // 不启用压缩的功能
		LocalTime:  true,      // 本地时间分割
	}
}

// parse log level
// default info level
func parseLogLevel(levelStr string) (logLevel zapcore.Level, err error) {
	// logLevel = zap.infoLevel
	var lvs = append(make([]string, 0), levelDebug, levelInfo, levelWarn, levelError, levelFatal)
	for _, v := range lvs {
		if !(levelStr == v) {
			continue
		}
		if logLevel, err = zapcore.ParseLevel(v); err == nil {
			return logLevel, nil
		}
	}

	return zap.InfoLevel, nil
}
