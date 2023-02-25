package zapimpl

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/llmuz/yggdrasill/ull"
)

type Option func(c *ZapHelperBuilder)

func AddHook(hook ull.Hook) Option {
	return func(c *ZapHelperBuilder) {
		c.hooks.Add(hook)
	}
}

type ZapHelperBuilder struct {
	hooks ull.LevelHooks
}

func NewHelper(logger *zap.Logger, opt ...Option) ull.Helper {
	var cfg = &ZapHelperBuilder{hooks: make(ull.LevelHooks, 6)}
	for _, v := range opt {
		v(cfg)
	}

	c := zapLoggerHelper{
		logger: logger,
		Hooks:  cfg.hooks,
	}
	c.initLogLevel()
	return &c
}

type zapLoggerHelper struct {
	logger *zap.Logger
	Hooks  ull.LevelHooks
	level  ull.Level
}

func (c *zapLoggerHelper) WithContext(ctx context.Context) ull.FieldLogger {
	return &zapFieldLogger{
		ctx:    ctx,
		helper: c,
		entry:  newZapLogEntry(ctx),
	}
}

func (c *zapLoggerHelper) levelEnabled(level ull.Level) bool {
	return c.level <= level
}

// 初始化 level 值
func (c *zapLoggerHelper) initLogLevel() {
	var levels = []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
	}

	for _, v := range levels {
		if c.logger.Core().Enabled(v) {
			c.level = ull.Level(v)
			break
		}
	}
}

func newZapLogEntry(ctx context.Context) ull.Entry {
	return &zapLoggerEntry{ctx: ctx, fields: make([]ull.Field, 0, 4)}
}
