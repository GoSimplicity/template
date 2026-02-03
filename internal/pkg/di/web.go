package di

import (
	"github.com/GoSimplicity/template/internal/api/http"
	"github.com/gin-gonic/gin"
)

// InitWeb 初始化web服务
func InitWeb(
	m []gin.HandlerFunc,
	templateHandler *http.TemplateHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(m...)

	templateHandler.RegisterRoutes(server)

	return server
}
