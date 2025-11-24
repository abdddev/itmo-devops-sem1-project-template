package v0

import "net/http"

type PricesAPI interface {
	UploadPrices(http.ResponseWriter, *http.Request)
	DownloadPrices(http.ResponseWriter, *http.Request)
}
