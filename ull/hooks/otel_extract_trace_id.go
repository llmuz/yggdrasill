package hooks

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/llmuz/yggdrasill/ull"
)

func NewOtelUllHook(lvs []ull.Level) (h ull.Hook) {
	return &otelUllHook{lvs: lvs}
}

type otelUllHook struct {
	lvs []ull.Level
}

func (c *otelUllHook) Levels() (lvs []ull.Level) {
	return c.lvs
}

func (c *otelUllHook) Fire(e ull.Entry) (err error) {
	spanCtx := trace.SpanContextFromContext(e.Context())
	if !spanCtx.IsValid() {
		return err
	}
	// 从 SpanContext 中提取 Trace ID 和 Span ID
	if err = e.AppendField(ull.Any("trace_id", spanCtx.TraceID())); err != nil {
		return err
	}

	if err = e.AppendField(ull.Any("span_id", spanCtx.SpanID())); err != nil {
		return err
	}
	return err
}
