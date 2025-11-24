package v0

import (
	"project_sem/internal/service"
)

type api struct {
	pricesService service.PricesService
}

func NewAPI(pricesService service.PricesService) *api {
	return &api{
		pricesService: pricesService,
	}
}
