package di

import (
	"github.com/GoSimplicity/template/internal/event"
	"github.com/gin-gonic/gin"
)

type Cmd struct {
	Server   *gin.Engine
	Consumer []event.Consumer
}
