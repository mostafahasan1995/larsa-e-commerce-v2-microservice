package config_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"github.com/samber/do"
)

type SettingRepo interface {
	Upsert() *helpers.ApiError
	FindOneAndIncreaseSerialNumber(ctx context.Context, collectionName string) (*config_models.SettingResponse, *helpers.ApiError)
	TruncateCollections(databaseName string, dto config_dtos.TruncateCollectionDtoList) *helpers.ApiError
	GetRelevantDbs() ([]string, *helpers.ApiError)
	AddContext(db string) context.Context
}

type SettingService interface {
	Upsert() *helpers.ApiError
	FindOneAndIncreaseSerialNumber(ctx context.Context, collectionName string) (*config_models.SettingResponse, *helpers.ApiError)
	TruncateCollections(databaseName string, dto config_dtos.TruncateCollectionDtoList) *helpers.ApiError
	GetRelevantDbs() ([]string, *helpers.ApiError)
	AddContext(db string) context.Context
}

type setting struct {
	repo SettingRepo
}

func NewSettingService(i *do.Injector) (SettingService, error) {
	return &setting{
		repo: do.MustInvoke[SettingRepo](i),
	}, nil
}

func (p *setting) Upsert() *helpers.ApiError {
	return p.repo.Upsert()
}

func (p *setting) FindOneAndIncreaseSerialNumber(ctx context.Context, collectionName string) (*config_models.SettingResponse, *helpers.ApiError) {
	return p.repo.FindOneAndIncreaseSerialNumber(ctx, collectionName)
}

func (p *setting) GetRelevantDbs() ([]string, *helpers.ApiError) {
	return p.repo.GetRelevantDbs()
}

func (p *setting) AddContext(db string) context.Context {
	return p.repo.AddContext(db)
}

func (p *setting) TruncateCollections(databaseName string, dto config_dtos.TruncateCollectionDtoList) *helpers.ApiError {
	return p.repo.TruncateCollections(databaseName, dto)
}
