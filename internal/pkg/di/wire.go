//go:build wireinject

package di

import (
	"github.com/GoSimplicity/template/internal/api/http"
	eventTemplate "github.com/GoSimplicity/template/internal/event/template"
	"github.com/GoSimplicity/template/internal/repository"
	"github.com/GoSimplicity/template/internal/repository/cache"
	"github.com/GoSimplicity/template/internal/repository/dao"
	"github.com/GoSimplicity/template/internal/service"
	"github.com/google/wire"
	_ "github.com/google/wire"
)

func InitializeApp() *Cmd {
	wire.Build(
		InitLogger,
		InitDB,
		InitRedis,
		InitSaramaClient,
		InitSyncProducer,
		InitConsumers,
		InitWeb,
		InitMiddlewares,
		http.NewTemplateHandler,
		service.NewTemplateService,
		repository.NewTemplateRepository,
		dao.NewTemplateDAO,
		cache.NewTemplateCache,
		eventTemplate.NewTemplateSaramaSyncProducer,
		eventTemplate.NewTemplateEventConsumer,
		eventTemplate.NewTemplateDeadLetterConsumer,
		wire.Struct(new(Cmd), "*"),
	)
	return new(Cmd)
}
