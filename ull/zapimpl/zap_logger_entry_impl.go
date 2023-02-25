package zapimpl

import (
	"context"

	"github.com/llmuz/yggdrasill/ull"
)

type zapLoggerEntry struct {
	ctx    context.Context // context
	fields []ull.Field     // data gen by Entry impl
}

// Context get ctx
func (c *zapLoggerEntry) Context() (ctx context.Context) {
	return c.ctx
}

// AppendField append field to fields
func (c *zapLoggerEntry) AppendField(field ull.Field) (err error) {
	c.fields = append(c.fields, field)
	return nil
}

// GetFields get fields
func (c *zapLoggerEntry) GetFields() []ull.Field {
	return c.fields
}
