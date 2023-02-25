package tracing

import (
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/logger"

	"github.com/llmuz/yggdrasill/ugorm"
	"github.com/llmuz/yggdrasill/ull"
)

func NewOtelHook() ugorm.Hook {
	return &otelTracingHook{levels: append(make([]logger.LogLevel, 0), logger.Info, logger.Warn, logger.Error)}
}

type otelTracingHook struct {
	levels []logger.LogLevel
}

func (c *otelTracingHook) Fire(e ull.Entry) (err error) {
	if cv := trace.SpanContextFromContext(e.Context()); cv.IsValid() {
		err = e.AppendField(ull.Any("trace_id", cv.TraceID().String()))
	}
	return err
}

func (c *otelTracingHook) Levels() (levels []logger.LogLevel) {
	return c.levels
}
