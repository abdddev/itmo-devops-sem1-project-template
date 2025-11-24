package archive

import "fmt"

const (
	TypeZip = "zip"
	TypeTar = "tar"
)

func NewHandler(archiveType string) (Handler, error) {
	switch archiveType {
	case TypeZip, "":
		return NewZipHandler(), nil
	case TypeTar:
		return NewTarHandler(), nil
	default:
		return nil, fmt.Errorf("unsupported archive type: %s", archiveType)
	}
}
