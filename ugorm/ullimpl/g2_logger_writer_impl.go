package ullimpl

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"

	"github.com/llmuz/yggdrasill/ugorm"
	"github.com/llmuz/yggdrasill/ull"
	"github.com/llmuz/yggdrasill/ull/zapimpl"
)

func NewWriter(logger *zap.Logger, optHooks ...ugorm.Hook) (w Writer) {
	var hooks = make(ugorm.Hooks, 0)
	for _, h := range optHooks {
		for _, level := range h.Levels() {
			hooks[level] = append(hooks[level], h)
		}
	}
	return &writerImpl{logger: logger, hooks: hooks}
}

type writerImpl struct {
	logger *zap.Logger
	hooks  ugorm.Hooks
}

func (c *writerImpl) Printf(ctx context.Context, level logger.LogLevel, format string, values ...interface{}) {
	var e = zapimpl.NewZapLogEntry(ctx)
	c.hooks.Fire(level, e)
	var f = make([]zapcore.Field, 0, 8)
	var fields = make([]ull.Field, 0)
	for _, v := range append(fields, e.GetFields()...) {
		f = append(f, zap.Any(v.Key, v.Interface))
	}
	var buf = strings.Builder{}
	buf.WriteString(fmt.Sprintf(format, values...))
	logLevel := checkLogLevel(level)
	if ce := c.logger.Check(logLevel, buf.String()); ce != nil {
		ce.Write(f...)
	}
}

func checkLogLevel(level logger.LogLevel) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.WarnLevel
	case Silent:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}
