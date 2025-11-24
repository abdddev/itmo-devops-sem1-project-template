package v0

import (
	"io"
	"log"
	"net/http"
)

func (a *api) DownloadPrices(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")
	minPrice := r.URL.Query().Get("min")
	maxPrice := r.URL.Query().Get("max")

	archiveReader, err := a.pricesService.ProcessDownload(r.Context(), startDate, endDate, minPrice, maxPrice)
	if err != nil {
		log.Printf("failed to process download: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"data.zip\"")
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, archiveReader); err != nil {
		log.Printf("failed to write archive to response: %v", err)
	}
}
