package prices

import (
	"context"
	"fmt"
	"io"

	"project_sem/internal/archive"
	def "project_sem/internal/service"
)

func (s *service) ProcessUpload(ctx context.Context, archiveData io.Reader, archiveType string) (def.UploadResult, error) {
	handler, err := archive.NewHandler(archiveType)
	if err != nil {
		return def.UploadResult{}, fmt.Errorf("failed to create archive handler: %w", err)
	}

	prices, totalCount, err := handler.Extract(archiveData)
	if err != nil {
		return def.UploadResult{}, fmt.Errorf("failed to extract archive: %w", err)
	}

	result, err := s.pricesRepository.InsertBatchTx(ctx, prices)
	if err != nil {
		return def.UploadResult{}, fmt.Errorf("failed to insert batch: %w", err)
	}

	return def.UploadResult{
		TotalCount:      totalCount,
		DuplicatesCount: result.DuplicatesCount,
		Summary:         result.Summary,
	}, nil
}
