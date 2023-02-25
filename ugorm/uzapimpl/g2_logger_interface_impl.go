package uzapimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Colors
const (
	Reset       = logger.Reset
	Red         = logger.Red
	Green       = logger.Green
	Yellow      = logger.Yellow
	Blue        = logger.Blue
	Magenta     = logger.Magenta
	Cyan        = logger.Cyan
	White       = logger.White
	BlueBold    = logger.BlueBold
	MagentaBold = logger.MagentaBold
	RedBold     = logger.RedBold
	YellowBold  = logger.YellowBold
)

const (
	// Silent silent log level
	Silent = logger.Silent
	// Info info log level
	Info = logger.Info
	// Warn warn log level
	Warn = logger.Warn
	// Error error log level
	Error = logger.Error
)

// ParserLevel 从字符串中解析 gorm 的日志等级
// 默认值为 info
func ParserLevel(lev string) (level logger.LogLevel) {
	lev = strings.ToLower(lev)
	var mapper = make(map[string]logger.LogLevel, 4)
	mapper["silent"] = Silent
	mapper["info"] = Info
	mapper["warn"] = Warn
	mapper["error"] = Error

	if v, ok := mapper[lev]; ok {
		level = v
	} else {
		level = Info
	}
	return level
}

func NewLoggerInterface(writer Writer, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s [info] "
		warnStr      = "%s [warn] "
		errStr       = "%s [error] "
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s " + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s " + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s " + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s " + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s " + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &g2LoggerInterfaceImpl{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type g2LoggerInterfaceImpl struct {
	logger.Config
	Writer
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

type Writer interface {
	Printf(ctx context.Context, level logger.LogLevel, format string, values ...interface{})
}

func (c *g2LoggerInterfaceImpl) LogMode(level logger.LogLevel) logger.Interface {
	keylogger := *c
	keylogger.LogLevel = level
	return &keylogger
}

func (c *g2LoggerInterfaceImpl) Info(ctx context.Context, format string, data ...interface{}) {
	if c.LogLevel >= Info {
		c.Printf(ctx, Info, c.infoStr+format, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (c *g2LoggerInterfaceImpl) Warn(ctx context.Context, format string, data ...interface{}) {
	if c.LogLevel >= Warn {
		c.Printf(ctx, Warn, c.warnStr+format, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (c *g2LoggerInterfaceImpl) Error(ctx context.Context, format string, data ...interface{}) {
	if c.LogLevel >= Error {
		c.Printf(ctx, Error, c.errStr+format, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (c *g2LoggerInterfaceImpl) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if c.LogLevel <= Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && c.LogLevel >= Error && (!errors.Is(err, logger.ErrRecordNotFound) || !c.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			c.Printf(ctx, Info, c.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			c.Printf(ctx, Info, c.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", c.SlowThreshold)
		if rows == -1 {
			c.Printf(ctx, Info, c.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			c.Printf(ctx, Info, c.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case c.LogLevel == Info:
		sql, rows := fc()
		if rows == -1 {
			c.Printf(ctx, Info, c.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			c.Printf(ctx, Info, c.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
