package prices

import (
	"context"
	"project_sem/internal/model"
)

func (r *repository) InsertBatch(ctx context.Context, prices []model.Price) (inserted int64, duplicated int64, err error) {
	if len(prices) == 0 {
		return 0, 0, nil
	}

	seenIDs := make(map[int64]bool)
	var uniquePrices []model.Price
	inputDuplicates := int64(0)

	for _, price := range prices {
		if seenIDs[price.ID] {
			inputDuplicates++
			continue
		}
		seenIDs[price.ID] = true
		uniquePrices = append(uniquePrices, price)
	}

	tx, err := r.db.Conn().BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO prices (id, name, category, price, create_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	dbDuplicates := int64(0)
	for _, price := range uniquePrices {
		result, err := stmt.ExecContext(ctx,
			price.ID,
			price.Name,
			price.Category,
			price.Price,
			price.CreateDate,
		)
		if err != nil {
			return 0, 0, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, 0, err
		}

		if rowsAffected == 0 {
			dbDuplicates++
		} else {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	duplicated = inputDuplicates + dbDuplicates
	return inserted, duplicated, nil
}
