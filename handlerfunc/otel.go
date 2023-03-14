package handlerfunc

import (
	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "github.com/llmuz/yggdrasill/handlerfunc"
)

// WithTraceProvider 参考文档
// https://github.com/open-telemetry/opentelemetry-go/edit/main/example/jaeger/main.go
func WithTraceProvider(tp *tracesdk.TracerProvider) gin.HandlerFunc {
	//tp := tracesdk.NewTracerProvider(tracesdk.WithSampler(tracesdk.NeverSample()))
	return func(ctx *gin.Context) {
		// 如果已经有 trace, 那就不用注入
		if v := trace.SpanFromContext(ctx.Request.Context()).SpanContext(); v.IsValid() {
			ctx.Header("trace_id", v.TraceID().String())
			return
		}

		newCtx, span := tp.Tracer(defaultTracerName).Start(ctx.Request.Context(), "pipeline-handler")
		defer span.End()
		ctx.Request = ctx.Request.WithContext(newCtx)
		if v := trace.SpanFromContext(ctx.Request.Context()).SpanContext(); v.IsValid() {
			ctx.Header("trace_id", v.TraceID().String())
		}

		ctx.Next()
	}
}
