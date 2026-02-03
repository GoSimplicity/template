package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/GoSimplicity/template/internal/domain"
	templateEvent "github.com/GoSimplicity/template/internal/event/template"
	"github.com/GoSimplicity/template/internal/repository/cache"
	"github.com/GoSimplicity/template/internal/repository/dao"
	"go.uber.org/zap"
)

type TemplateRepository interface {
	Create(ctx context.Context, template domain.Template) (domain.Template, error)
}

type templateRepository struct {
	dao      dao.TemplateDAO
	cache    cache.TemplateCache
	logger   *zap.Logger
	producer templateEvent.Producer
}

func NewTemplateRepository(dao dao.TemplateDAO, logger *zap.Logger, cache cache.TemplateCache, producer templateEvent.Producer) TemplateRepository {
	return &templateRepository{
		dao:      dao,
		logger:   logger,
		cache:    cache,
		producer: producer,
	}
}

func (r *templateRepository) Create(ctx context.Context, template domain.Template) (domain.Template, error) {
	created, err := r.dao.Create(ctx, toEntity(template))
	if err != nil {
		r.logger.Error("create template failed", zap.Error(err))
		r.producer.ProduceTemplateEvent(templateEvent.TemplateEvent{
			TemplateId: template.ID,
		})
		return domain.Template{}, err
	}

	r.cache.Set(ctx, fmt.Sprintf("template:%d", created.ID), strconv.FormatInt(created.ID, 10), time.Hour*24)
	return toDomain(created), nil
}

func toDomain(daoTemplate dao.Template) domain.Template {
	return domain.Template{
		ID: daoTemplate.ID,
	}
}

func toEntity(domainTemplate domain.Template) dao.Template {
	return dao.Template{
		ID: domainTemplate.ID,
	}
}
