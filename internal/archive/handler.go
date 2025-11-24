package archive

import (
	"io"

	"project_sem/internal/model"
)

type Handler interface {
	Extract(r io.Reader) (prices []model.Price, totalCount int64, err error)
	Create(prices []model.Price) (io.Reader, error)
}
