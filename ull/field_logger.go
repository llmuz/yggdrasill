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
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}
