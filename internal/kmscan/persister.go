package kmscan

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/maicher/kmscan/internal/monitor"
)

type Persister interface {
	Persist(image.Image, string) error
}

type FilePersister struct {
	dirPath string
	monitor *monitor.Monitor
}

func NewFilePersister(dirPath string, m *monitor.Monitor) (*FilePersister, error) {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("Error creating directory: %s", err)
	}

	return &FilePersister{dirPath: dirPath, monitor: m}, nil
}

func (p FilePersister) Persist(img image.Image, name string) error {
	t := time.Now()

	outputPath := filepath.Join(p.dirPath, name)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %s", err)
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, img, nil); err != nil {
		return fmt.Errorf("Error saving image: %s", err)
	}

	p.monitor.Processor(time.Since(t), "%s saved", name)

	return nil
}

type NullPersister struct {
}

func (p NullPersister) Persist(img image.Image, name string) error {
	return nil
}
