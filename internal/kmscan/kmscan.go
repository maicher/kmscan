package kmscan

import (
	"context"
	"image"
	"os"
	"time"

	"github.com/maicher/kmscan/internal/ui"
)

type Kmscan struct {
	Scanner        *Scanner
	Filters        *Filters
	Extractor      *Extractor
	PhotoPersister Persister
	ScanPersister  Persister
	Autorotator    *Autorotator
	Uploader       Uploader
	Keyboard       *ui.Keyboard
	Logger         *ui.Logger
}

func (k *Kmscan) Run(devices []string) {
	scansCh := make(chan string, 10)
	photosCh := make(chan Photo, 10)
	ctx, cancel := context.WithCancel(context.Background())

	// Perform scanning and write results to scansCh.
	for i, device := range devices {
		go k.loopScanning(ctx, device, i+1, scansCh)
	}

	// Read from scansCh, extract individual photos and write results to photosCh.
	go k.loopScansProcessing(scansCh, photosCh)

	// Read from photosCh, process the photos and persist them..
	go k.loopPhotosProcessing(photosCh)

	// Wait keyboard commands.
	if err := k.Keyboard.Listen(); err != nil {
		k.Logger.Err("unable to initialize keyboard: %s", err)
	}

	k.Logger.Msg("", "quitting...")
	cancel()
	time.Sleep(1 * time.Second)
}

func (k *Kmscan) ProcessImage(imagePath string) {
	photosCh := make(chan Photo, 10)

	scan, err := NewScan(imagePath)
	if err != nil {
		k.Logger.Err("%s", err)
		return
	}

	k.Logger.Msg(imagePath, "file loaded")

	k.Filters.ApplyFilters(scan)
	if err := k.persistScans(scan); err != nil {
		k.Logger.Err("%s", err)
	}

	go func() {
		k.Extractor.Extract(scan, photosCh)
		close(photosCh)
	}()

	k.loopPhotosProcessing(photosCh)
	time.Sleep(1 * time.Second)

}

func (k *Kmscan) loopScanning(ctx context.Context, device string, deviceNo int, scansCh chan<- string) {
	var (
		i        int
		scanPath string
		err      error
	)

	for {
		scanPath, err = k.Scanner.Scan(ctx, device, deviceNo, i)
		if err != nil {
			break
		}

		scansCh <- scanPath
		i++
		time.Sleep(100 * time.Millisecond)
	}
}

func (k *Kmscan) loopScansProcessing(scansCh chan string, photosCh chan Photo) {
	for imagePath := range scansCh {
		scan, err := NewScan(imagePath)
		if err != nil {
			k.Logger.Err("%s", err)
			continue
		}
		defer os.Remove(imagePath)

		k.Filters.ApplyFilters(scan)
		if err := k.persistScans(scan); err != nil {
			k.Logger.Err("%s", err)
		}

		k.Extractor.Extract(scan, photosCh)
	}
}

func (k *Kmscan) loopPhotosProcessing(photosCh <-chan Photo) {
	var img image.Image

	for photo := range photosCh {
		img = k.Autorotator.Autorotate(photo.Image, photo.Name)

		k.PhotoPersister.Persist(img, photo.Name)

		if err := k.Uploader.Upload(photo.Name); err != nil {
			time.Sleep(100 * time.Millisecond)
			k.Uploader.Upload(photo.Name)
		}
	}
}

func (k *Kmscan) persistScans(scan *Scan) error {
	if err := k.ScanPersister.Persist(scan.Image, "_0_scan.jpg"); err != nil {
		return err
	}

	if err := k.ScanPersister.Persist(scan.Maximum, "_1_max.jpg"); err != nil {
		return err
	}

	if err := k.ScanPersister.Persist(scan.Gray, "_2_gray.jpg"); err != nil {
		return err
	}

	if err := k.ScanPersister.Persist(scan.Binary, "_3_binary.jpg"); err != nil {
		return err
	}

	return nil
}

func maxIndex(arr [4]int) int {
	maxIdx := 0
	maxValue := arr[0]

	for i, value := range arr {
		if value > maxValue {
			maxValue = value
			maxIdx = i
		}
	}

	return maxIdx
}
