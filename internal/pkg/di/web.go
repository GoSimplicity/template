package di

import (
	"github.com/GoSimplicity/template/internal/api/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitWeb 初始化web服务
func InitWeb(
	m []gin.HandlerFunc,
	templateHandler *http.TemplateHandler,
) *gin.Engine {
	setGinModeFromConfig()

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(m...)

	templateHandler.RegisterRoutes(server)

	return server
}

func setGinModeFromConfig() {
	mode := viper.GetString("server.mode")
	switch mode {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		gin.SetMode(mode)
	default:
		if mode == "" {
			gin.SetMode(gin.DebugMode)
			return
		}
		gin.SetMode(gin.DebugMode)
	}
}
