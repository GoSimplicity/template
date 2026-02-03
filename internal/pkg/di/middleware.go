package di

import (
	"github.com/GoSimplicity/template/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InitMiddlewares 初始化中间件
func InitMiddlewares(logger *zap.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.CORS(),
		middleware.RequestLogger(logger),
	}
}
