package v0

import (
	"encoding/json"
	"log"
	"net/http"
)

type UploadPricesResponse struct {
	TotalCount      int64   `json:"total_count"`
	DuplicatesCount int64   `json:"duplicates_count"`
	TotalItems      int64   `json:"total_items"`
	TotalCategories int64   `json:"total_categories"`
	TotalPrice      float64 `json:"total_price"`
}

func (a *api) UploadPrices(w http.ResponseWriter, r *http.Request) {
	archiveType := r.URL.Query().Get("type")
	if archiveType == "" {
		archiveType = "zip"
	}

	if archiveType != "zip" && archiveType != "tar" {
		http.Error(w, "invalid archive type: must be 'zip' or 'tar'", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("failed to parse multipart form: %v", err)
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("failed to get file from form: %v", err)
		http.Error(w, "file not found in request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	const maxFileSize = 32 << 20
	if header.Size > maxFileSize {
		log.Printf("file too large: %d bytes (max: %d bytes)", header.Size, maxFileSize)
		http.Error(w, "file too large: maximum size is 30 MB", http.StatusRequestEntityTooLarge)
		return
	}

	result, err := a.pricesService.ProcessUpload(r.Context(), file, archiveType)
	if err != nil {
		log.Printf("failed to process upload: %v", err)
		http.Error(w, "failed to process upload", http.StatusInternalServerError)
		return
	}

	resp := UploadPricesResponse{
		TotalCount:      result.TotalCount,
		DuplicatesCount: result.DuplicatesCount,
		TotalItems:      result.Summary.TotalItems,
		TotalCategories: result.Summary.TotalCategories,
		TotalPrice:      result.Summary.TotalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
