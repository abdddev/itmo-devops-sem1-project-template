package prices

import (
	"context"
	"project_sem/internal/model"
)

func (r *repository) InsertBatchTx(ctx context.Context, prices []model.Price) (result model.InsertResult, err error) {
	tx, err := r.db.Conn().BeginTx(ctx, nil)
	if err != nil {
		return model.InsertResult{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var dupCount int64

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO prices (name, category, price, create_date)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (name, category, price, create_date) DO NOTHING
    `)
	if err != nil {
		return model.InsertResult{}, err
	}
	defer stmt.Close()

	for _, p := range prices {
		res, errExec := stmt.ExecContext(ctx, p.Name, p.Category, p.Price, p.CreateDate)
		if errExec != nil {
			err = errExec
			return model.InsertResult{}, err
		}

		affected, errRows := res.RowsAffected()
		if errRows != nil {
			err = errRows
			return model.InsertResult{}, err
		}
		if affected == 0 {
			dupCount++
		}
	}

	var summary model.PriceSummary
	err = tx.QueryRowContext(ctx, `
        SELECT 
            COUNT(*)                AS total_items,
            COUNT(DISTINCT category) AS total_categories,
            COALESCE(SUM(price), 0) AS total_price
        FROM prices
    `).Scan(&summary.TotalItems, &summary.TotalCategories, &summary.TotalPrice)
	if err != nil {
		return model.InsertResult{}, err
	}

	if err = tx.Commit(); err != nil {
		return model.InsertResult{}, err
	}

	return model.InsertResult{
		DuplicatesCount: dupCount,
		Summary:         summary,
	}, nil
}
