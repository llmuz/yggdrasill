package uzapimpl

import (
	"context"

	"github.com/llmuz/yggdrasill/ull"
)

type entryImpl struct {
	ctx    context.Context
	fields []ull.Field
}

func (c *entryImpl) Context() (ctx context.Context) {
	return c.ctx
}

func (c *entryImpl) AppendField(field ull.Field) (err error) {
	c.fields = append(c.fields, field)
	return nil
}

func (c *entryImpl) GetFields() (fields []ull.Field) {
	return c.fields
}

func newEntry(ctx context.Context) (f ull.Entry) {
	return &entryImpl{
		ctx:    ctx,
		fields: make([]ull.Field, 0),
	}
}
