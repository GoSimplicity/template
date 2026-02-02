package service

import (
	"context"
	"go_project_template/internal/domain"
	"go_project_template/internal/repository"
)

type TemplateService interface {
	CreateTemplate(ctx context.Context, req domain.Template) (domain.Template, error)
}

type templateService struct {
	repo repository.TemplateRepository
}

func NewTemplateService(repo repository.TemplateRepository) TemplateService {
	return &templateService{
		repo: repo,
	}
}

func (s *templateService) CreateTemplate(ctx context.Context, req domain.Template) (domain.Template, error) {
	created, err := s.repo.Create(ctx, req)
	if err != nil {
		return domain.Template{}, err
	}
	req.ID = created.ID
	return created, nil
}
