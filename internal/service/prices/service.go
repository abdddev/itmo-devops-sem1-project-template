package prices

import (
	"project_sem/internal/repository"
	def "project_sem/internal/service"
)

var _ def.PricesService = (*service)(nil)

type service struct {
	pricesRepository repository.PricesRepository
}

func NewService(pricesRepository repository.PricesRepository) *service {
	return &service{
		pricesRepository: pricesRepository,
	}
}
