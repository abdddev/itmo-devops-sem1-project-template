package prices

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"project_sem/internal/archive"
	"project_sem/internal/model"
)

func (s *service) ProcessDownload(ctx context.Context, startDate, endDate, minPrice, maxPrice string) (io.Reader, error) {
	var filter model.Filter

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format, expected YYYY-MM-DD: %w", err)
		}
		filter.StartDate = start
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format, expected YYYY-MM-DD: %w", err)
		}
		filter.EndDate = end
	}

	if minPrice != "" {
		min, err := strconv.ParseInt(minPrice, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min price format: %w", err)
		}
		if min <= 0 {
			return nil, fmt.Errorf("min price must be positive")
		}
		filter.MinPrice = min
	}

	if maxPrice != "" {
		max, err := strconv.ParseInt(maxPrice, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max price format: %w", err)
		}
		if max <= 0 {
			return nil, fmt.Errorf("max price must be positive")
		}
		filter.MaxPrice = max
	}

	prices, err := s.pricesRepository.FindByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find prices: %w", err)
	}

	zipHandler := archive.NewZipHandler()
	archiveReader, err := zipHandler.Create(prices)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}

	return archiveReader, nil
}
