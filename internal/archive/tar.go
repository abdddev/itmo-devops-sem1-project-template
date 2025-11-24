package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"

	"project_sem/internal/model"
)

type TarHandler struct{}

func NewTarHandler() *TarHandler {
	return &TarHandler{}
}

func (h *TarHandler) Extract(r io.Reader) ([]model.Price, int64, error) {
	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("failed to read tar header: %w", err)
		}

		if header.Typeflag == tar.TypeDir {
			continue
		}

		if len(header.Name) > 4 && header.Name[len(header.Name)-4:] == ".csv" {
			return parseCSV(tarReader)
		}
	}

	return nil, 0, fmt.Errorf("CSV file not found in archive")
}

func (h *TarHandler) Create(prices []model.Price) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)

	csvBuf := new(bytes.Buffer)
	if err := writeCSV(csvBuf, prices); err != nil {
		return nil, err
	}

	header := &tar.Header{
		Name: "data.csv",
		Mode: 0644,
		Size: int64(csvBuf.Len()),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := tarWriter.Write(csvBuf.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write csv content: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}
