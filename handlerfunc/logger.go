package handlerfunc

import (
	"fmt"
	"strings"
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

		sb := strings.Builder{}
		for _, v := range append(defaultLogFormatter(param), biz(c)...) {
			sb.WriteString(fmt.Sprintf("%s:%+v ", v.Key, v.Interface))
		}
		logger.WithContext(c.Request.Context()).Infof(sb.String())
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

var defaultLogFormatter = func(param gin.LogFormatterParams) (fields []ull.Field) {

	fields = make([]ull.Field, 0, 8)
	fields = append(fields, ull.Any("status_code", param.StatusCode))
	fields = append(fields, ull.Any("latency_seconds", param.Latency.Seconds()))
	fields = append(fields, ull.Any("client_ip", param.ClientIP))
	fields = append(fields, ull.Any("method", param.Method))
	fields = append(fields, ull.Any("path", param.Path))
	fields = append(fields, ull.Any("body_size", param.BodySize))
	fields = append(fields, ull.Any("error_message", param.ErrorMessage))
	return fields

}
