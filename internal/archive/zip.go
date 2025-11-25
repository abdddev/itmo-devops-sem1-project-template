package archive

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"project_sem/internal/model"
)

type ZipHandler struct{}

func NewZipHandler() *ZipHandler {
	return &ZipHandler{}
}

func (h *ZipHandler) Extract(r io.Reader) ([]model.Price, int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read archive: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create zip reader: %w", err)
	}

	var csvFile *zip.File
	for _, f := range zipReader.File {
		if len(f.Name) > 4 && f.Name[len(f.Name)-4:] == ".csv" {
			csvFile = f
			break
		}
	}

	if csvFile == nil {
		return nil, 0, fmt.Errorf("CSV file not found in archive")
	}

	rc, err := csvFile.Open()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open data.csv: %w", err)
	}
	defer rc.Close()

	return parseCSV(rc)
}

func (h *ZipHandler) Create(prices []model.Price) (io.Reader, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	csvFile, err := zipWriter.Create("data.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to create data.csv in archive: %w", err)
	}

	if err := writeCSV(csvFile, prices); err != nil {
		zipWriter.Close()
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func parseCSV(r io.Reader) ([]model.Price, int64, error) {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read csv: %w", err)
	}

	var prices []model.Price
	var totalCount int64

	for i, record := range records {
		if i == 0 && record[0] == "id" {
			continue
		}

		totalCount++

		if len(record) != 5 {
			continue
		}

		id, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		name := record[1]
		category := record[2]
		if name == "" || category == "" {
			continue
		}

		price, err := strconv.ParseFloat(record[3], 64)
		if err != nil || price <= 0 {
			continue
		}

		createDate, err := time.Parse("2006-01-02", record[4])
		if err != nil {
			continue
		}

		prices = append(prices, model.Price{
			ID:         id,
			Name:       name,
			Category:   category,
			Price:      price,
			CreateDate: createDate,
		})
	}

	return prices, totalCount, nil
}

func writeCSV(w io.Writer, prices []model.Price) error {
	csvWriter := csv.NewWriter(w)

	if err := csvWriter.Write([]string{"id", "name", "category", "price", "create_date"}); err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}

	for _, p := range prices {
		record := []string{
			strconv.FormatInt(p.ID, 10),
			p.Name,
			p.Category,
			strconv.FormatFloat(p.Price, 'f', 2, 64),
			p.CreateDate.Format("2006-01-02"),
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write csv record: %w", err)
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
