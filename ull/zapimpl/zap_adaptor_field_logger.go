package zapimpl

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/llmuz/yggdrasill/ull"
)

const (
	debugLevel = zapcore.DebugLevel
	infoLevel  = zapcore.InfoLevel
	warnLevel  = zapcore.WarnLevel
	errorLevel = zapcore.ErrorLevel
	fatalLevel = zapcore.FatalLevel
)

type zapFieldLogger struct {
	entry  ull.Entry
	helper *zapLoggerHelper
	ctx    context.Context
}

func (c *zapFieldLogger) Debug(msg string, fields ...ull.Field) {
	if c.helper.levelEnabled(ull.Level(zapcore.DebugLevel)) {
		c.write(debugLevel, msg, fields...)
	}
}

func (c *zapFieldLogger) Info(msg string, fields ...ull.Field) {
	if c.helper.levelEnabled(ull.Level(infoLevel)) {
		c.write(infoLevel, msg, fields...)
	}
}

func (c *zapFieldLogger) Warn(msg string, fields ...ull.Field) {
	if c.helper.levelEnabled(ull.Level(warnLevel)) {
		c.write(warnLevel, msg, fields...)
	}
}

func (c *zapFieldLogger) Error(msg string, fields ...ull.Field) {
	if c.helper.levelEnabled(ull.Level(errorLevel)) {
		c.write(errorLevel, msg, fields...)
	}
}

func (c *zapFieldLogger) Fatal(msg string, fields ...ull.Field) {
	if c.helper.levelEnabled(ull.Level(fatalLevel)) {
		c.write(fatalLevel, msg, fields...)
	}
}

// fireHooks execute hook, if get error print error
// msg to the StdOut
func (c *zapFieldLogger) fireHooks() {
	var tmpHooks = make(ull.LevelHooks, len(c.helper.Hooks))
	for k, v := range c.helper.Hooks {
		tmpHooks[k] = v
	}
	if err := tmpHooks.Fire(c.helper.level, c.entry); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

// 写日志
func (c *zapFieldLogger) write(level zapcore.Level, msg string, fields ...ull.Field) {
	c.fireHooks()
	var f = make([]zapcore.Field, 0, 8)
	for _, v := range append(fields, c.entry.GetFields()...) {
		f = append(f, zap.Any(v.Key, v.Interface))
	}
	if ce := c.helper.logger.Check(level, msg); ce != nil {
		ce.Write(f...)
	}
}
