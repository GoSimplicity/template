package dao

import (
	"context"

	"gorm.io/gorm"
)

type Template struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;comment:'模板ID'"`
	CreatedAt int64
	UpdatedAt int64
}

func (Template) TableName() string {
	return "template"
}

type TemplateDAO interface {
	Create(ctx context.Context, template Template) (Template, error)
}

type templateDAO struct {
	db *gorm.DB
}

func NewTemplateDAO(db *gorm.DB) TemplateDAO {
	return &templateDAO{db: db}
}

func (d *templateDAO) Create(ctx context.Context, template Template) (Template, error) {
	if err := d.db.WithContext(ctx).Create(&template).Error; err != nil {
		return Template{}, err
	}
	return template, nil
}
