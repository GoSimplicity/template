package repository

import (
	"context"
	"go_project_template/internal/domain"
	"go_project_template/internal/repository/dao"

	"go.uber.org/zap"
)

type TemplateRepository interface {
	Create(ctx context.Context, template domain.Template) (domain.Template, error)
}

type templateRepository struct {
	dao    dao.TemplateDAO
	logger *zap.Logger
}

func NewTemplateRepository(dao dao.TemplateDAO, logger *zap.Logger) TemplateRepository {
	return &templateRepository{dao: dao, logger: logger}
}

func (r *templateRepository) Create(ctx context.Context, template domain.Template) (domain.Template, error) {
	created, err := r.dao.Create(ctx, toEntity(template))
	if err != nil {
		r.logger.Error("create template failed", zap.Error(err))
		return domain.Template{}, err
	}
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
