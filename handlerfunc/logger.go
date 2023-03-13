package handlerfunc

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/llmuz/yggdrasill/ull"
)

func LoggerHandlerFunc(logger ull.Helper, biz func(ctx *gin.Context) (fields []ull.Field)) gin.HandlerFunc {

	// 初始化
	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			//isTerm:  isTerm,
			Keys: c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		logger.WithContext(c.Request.Context()).Info(defaultLogFormatter(param), biz(c)...)
	}
}

var DefaultLoggerTracer = func(c *gin.Context) (fields []ull.Field) {
	if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
		fields = append(fields,
			ull.Any("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()),
		)
	}
	return fields
}

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006:01:02 - 15:04:05.666666"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
