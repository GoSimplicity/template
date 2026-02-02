package http

import (
	"net/http"

	"github.com/GoSimplicity/template/internal/service"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	service service.TemplateService
}

func NewTemplateHandler(service service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		service: service,
	}
}

func (h *TemplateHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/templates", h.CreateTemplate)
	// ...
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Template created"})
}
