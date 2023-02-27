package uzapimpl

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	"github.com/llmuz/yggdrasill/ugorm"
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
	var e = newEntry(ctx)
	c.hooks.Fire(level, e)
	var f = make([]zap.Field, 0, len(e.GetFields()))
	for _, v := range e.GetFields() {
		f = append(f, zap.Any(v.Key, v.Interface))
	}
	c.logger.Info(fmt.Sprintf(format, values...), f...)
}
