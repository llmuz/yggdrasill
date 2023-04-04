package ullimpl

import (
	"context"
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/llmuz/yggdrasill/ugorm"
	"github.com/llmuz/yggdrasill/ull"
	"github.com/llmuz/yggdrasill/ull/zapimpl"
)

func NewWriter(log ull.Helper, optHooks ...ugorm.Hook) (w Writer) {
	var hooks = make(ugorm.Hooks, 0)
	for _, h := range optHooks {
		for _, level := range h.Levels() {
			hooks[level] = append(hooks[level], h)
		}
	}
	return &writerImpl{log: log, hooks: hooks}
}

type writerImpl struct {
	log   ull.Helper
	hooks ugorm.Hooks
}

func (c *writerImpl) Printf(ctx context.Context, level logger.LogLevel, format string, values ...interface{}) {
	var e = zapimpl.NewZapLogEntry(ctx)
	c.hooks.Fire(level, e)

	c.log.WithContext(ctx).Info(fmt.Sprintf(format, values...), e.GetFields()...)
}
