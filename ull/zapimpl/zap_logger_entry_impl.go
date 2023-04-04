package zapimpl

import (
	"context"

	"github.com/llmuz/yggdrasill/ull"
)

type ZapLoggerEntry struct {
	ctx    context.Context // context
	fields []ull.Field     // data gen by Entry impl
}

// Context get ctx
func (c *ZapLoggerEntry) Context() (ctx context.Context) {
	return c.ctx
}

// AppendField append field to fields
func (c *ZapLoggerEntry) AppendField(field ull.Field) (err error) {
	c.fields = append(c.fields, field)
	return nil
}

// GetFields get fields
func (c *ZapLoggerEntry) GetFields() []ull.Field {
	return c.fields
}

func NewZapLogEntry(ctx context.Context) ull.Entry {
	return &ZapLoggerEntry{ctx: ctx, fields: make([]ull.Field, 0, 4)}
}
