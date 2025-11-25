package repository

import (
	"context"
	"project_sem/internal/model"
)

type PricesRepository interface {
	InsertBatchTx(ctx context.Context, prices []model.Price) (model.InsertResult, error)
	GetSummary(ctx context.Context) (model.PriceSummary, error)
	FindByFilter(ctx context.Context, filter model.Filter) ([]model.Price, error)
}
