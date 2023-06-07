package ull

import (
	"context"
)

func Any(key string, value interface{}) Field {
	return Field{Key: key, Interface: value}
}

type Field struct {
	Key       string
	Interface interface{}
}

type Entry interface {
	Context() (ctx context.Context)
	AppendField(field Field) (err error)
	GetFields() (fields []Field)
}

type Helper interface {
	WithContext(ctx context.Context) (logger FieldLogger)
}

type FieldLogger interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
	Fatalf(format string, a ...any)
}
