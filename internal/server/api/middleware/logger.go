package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/MaxRadzey/gog/internal/logger"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// RequestLogger возвращает Gin middleware для логирования входящих запросов.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Log.Info("got incoming HTTP request",
			zap.String("URI", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", duration),
		)
	}
}

// ResponseLogger возвращает Gin middleware для логирования ответов (status, size).
func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := &loggingResponseWriter{
			ResponseWriter: c.Writer,
			responseData:   responseData,
		}
		c.Writer = lw
		c.Next()

		logger.Log.Info("response",
			zap.Int("status", lw.responseData.status),
			zap.Int("size", lw.responseData.size),
		)
	}
}
