package config_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"

	"github.com/samber/do"
)

type MigrationRepo interface {
	MigrateData(ctx context.Context, dto config_dtos.DataMigrationDto) *helpers.ApiError
}

type MigrationService interface {
	MigrateData(ctx context.Context, dto config_dtos.DataMigrationDto) *helpers.ApiError
}

type migration struct {
	repo MigrationRepo
}

func NewMigrationService(i *do.Injector) (MigrationService, error) {
	return &migration{
		repo: do.MustInvoke[MigrationRepo](i),
	}, nil
}

func (m *migration) MigrateData(ctx context.Context, dto config_dtos.DataMigrationDto) *helpers.ApiError {
	return m.repo.MigrateData(ctx, dto)
}
