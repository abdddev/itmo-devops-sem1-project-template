package prices

import (
	"context"
	"fmt"
	"project_sem/internal/model"
)

func (r *repository) FindByFilter(ctx context.Context, filter model.Filter) ([]model.Price, error) {
	query := `SELECT id, name, category, price, create_date FROM prices WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if !filter.StartDate.IsZero() {
		query += ` AND create_date >= $` + fmt.Sprint(argIdx)
		args = append(args, filter.StartDate)
		argIdx++
	}

	if !filter.EndDate.IsZero() {
		query += ` AND create_date <= $` + fmt.Sprint(argIdx)
		args = append(args, filter.EndDate)
		argIdx++
	}

	if filter.MinPrice > 0 {
		query += ` AND price >= $` + fmt.Sprint(argIdx)
		args = append(args, filter.MinPrice)
		argIdx++
	}

	if filter.MaxPrice > 0 {
		query += ` AND price <= $` + fmt.Sprint(argIdx)
		args = append(args, filter.MaxPrice)
		argIdx++
	}

	query += ` ORDER BY id`

	rows, err := r.db.Conn().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []model.Price
	for rows.Next() {
		var p model.Price
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.CreateDate); err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}
