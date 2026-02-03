package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxLoggedBodyBytes = 1 << 20 // 1 MiB

// RequestLogger 记录请求详情和响应结果
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		bodyBytes, truncated := readBodyWithLimit(c.Request.Body, maxLoggedBodyBytes)
		if bodyBytes != nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("raw_query", c.Request.URL.RawQuery),
			zap.Int("status", status),
			zap.Int("resp_size", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.Duration("latency", latency),
			zap.Any("headers", sanitizeHeaders(c.Request.Header)),
		}

		if bodyBytes != nil && len(bodyBytes) > 0 {
			fields = append(fields,
				zap.ByteString("body", bodyBytes),
				zap.Bool("body_truncated", truncated),
			)
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			logger.Error("http request", fields...)
		case status >= 400:
			logger.Warn("http request", fields...)
		default:
			logger.Info("http request", fields...)
		}
	}
}

func readBodyWithLimit(body io.ReadCloser, limit int64) ([]byte, bool) {
	if body == nil {
		return nil, false
	}
	defer body.Close()

	limited := io.LimitReader(body, limit+1)
	data, _ := io.ReadAll(limited)
	if int64(len(data)) > limit {
		return data[:limit], true
	}
	return data, false
}

func sanitizeHeaders(headers map[string][]string) map[string][]string {
	if headers == nil {
		return nil
	}
	safe := make(map[string][]string, len(headers))
	for k, v := range headers {
		switch {
		case strings.EqualFold(k, "Authorization"),
			strings.EqualFold(k, "Cookie"),
			strings.EqualFold(k, "X-Refresh-Token"):
			safe[k] = []string{"***"}
		default:
			safe[k] = v
		}
	}
	return safe
}
