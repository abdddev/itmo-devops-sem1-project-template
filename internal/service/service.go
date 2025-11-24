package service

import (
	"context"
	"io"

	"project_sem/internal/model"
)

type UploadResult struct {
	TotalCount      int64
	DuplicatesCount int64
	Summary         model.PriceSummary
}

type PricesService interface {
	ProcessUpload(ctx context.Context, archiveData io.Reader, archiveType string) (UploadResult, error)
	ProcessDownload(ctx context.Context, startDate, endDate, minPrice, maxPrice string) (io.Reader, error)
}
