package fs

import (
	"xl/document"
)

type FileInterface interface {
	Open() (*document.Document, error)
	Write(*document.Document) error
}
