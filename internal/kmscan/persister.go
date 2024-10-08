package kmscan

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/maicher/kmscan/internal/ui"
)

type Persister interface {
	Persist(image.Image, string) error
}

type FilePersister struct {
	dirPath string
	ui      *ui.Logger
}

func NewFilePersister(dirPath string, m *ui.Logger) (*FilePersister, error) {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("error creating directory: %s", err)
	}

	return &FilePersister{dirPath: dirPath, ui: m}, nil
}

func (p FilePersister) Persist(img image.Image, name string) error {
	t := time.Now()

	outputPath := filepath.Join(p.dirPath, name)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err)
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, img, nil); err != nil {
		return fmt.Errorf("error saving image: %s", err)
	}

	p.ui.MsgWithDuration(time.Since(t), "%s saved", name)

	return nil
}

type NullPersister struct {
}

func (p NullPersister) Persist(img image.Image, name string) error {
	return nil
}
