package repository

import (
	"context"
	"project_sem/internal/model"
)

type PricesRepository interface {
	InsertBatch(ctx context.Context, prices []model.Price) (inserted int64, duplicated int64, err error)
	GetSummary(ctx context.Context) (model.PriceSummary, error)
	FindByFilter(ctx context.Context, filter model.Filter) ([]model.Price, error)
}
