package prices

import (
	"context"
	"project_sem/internal/model"
)

func (r *repository) GetSummary(ctx context.Context) (model.PriceSummary, error) {
	query := `
		SELECT
			COUNT(*) as total_items,
			COUNT(DISTINCT category) as total_categories,
			COALESCE(SUM(price), 0) as total_price
		FROM prices
	`

	var summary model.PriceSummary
	err := r.db.Conn().QueryRowContext(ctx, query).Scan(
		&summary.TotalItems,
		&summary.TotalCategories,
		&summary.TotalPrice,
	)
	if err != nil {
		return model.PriceSummary{}, err
	}

	return summary, nil
}
