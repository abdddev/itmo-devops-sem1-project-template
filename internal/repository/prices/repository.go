package prices

import (
	"project_sem/internal/infrastructure/database"
	def "project_sem/internal/repository"
)

var _ def.PricesRepository = (*repository)(nil)

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	return &repository{
		db: db,
	}
}
