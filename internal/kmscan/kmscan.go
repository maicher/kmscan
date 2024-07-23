package kmscan

import (
	"context"
	"image"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/maicher/kmscan/internal/monitor"
)

type Kmscan struct {
	Scanner            *Scanner
	Filters            *Filters
	Extractor          *Extractor
	FaceDetector       *FaceDetector
	PhotoPersister     Persister
	ScanPersister      Persister
	Uploader           *Uploader
	KeyboardController *KeyboardController
	Monitor            *monitor.Monitor
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
	if err := k.KeyboardController.Listen(); err != nil {
		k.Monitor.Err("unable to initialize keyboard: %s", err)
	}

	k.Monitor.Msg("", "quitting...")
	cancel()
	time.Sleep(1 * time.Second)
}

func (k *Kmscan) ProcessImage(imagePath string) {
	photosCh := make(chan Photo, 10)

	scan, err := NewScan(imagePath)
	if err != nil {
		k.Monitor.Err("%s", err)
		return
	}

	k.Monitor.Msg(imagePath, "file loaded")

	k.Filters.ApplyFilters(scan)
	if err := k.persistScans(scan); err != nil {
		k.Monitor.Err("%s", err)
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
			k.Monitor.Err("%s", err)
			continue
		}
		defer os.Remove(imagePath)

		k.Filters.ApplyFilters(scan)
		if err := k.persistScans(scan); err != nil {
			k.Monitor.Err("%s", err)
		}

		k.Extractor.Extract(scan, photosCh)
	}
}

func (k *Kmscan) loopPhotosProcessing(photosCh <-chan Photo) {
	imgs := [4]image.Image{}
	grays := [4]*image.Gray{}
	faces := [4]int{}
	var i int

	for photo := range photosCh {
		t := time.Now()
		imgs[0] = photo.Image
		grays[0] = k.Filters.Gray(imgs[0])
		grays[0] = k.Filters.Sharpen(grays[0])

		faces[0] = k.FaceDetector.Faces(grays[0])

		imgs[1] = k.Filters.Rotate(imgs[0])
		grays[1] = k.Filters.Gray(imgs[1])
		faces[1] = k.FaceDetector.Faces(grays[1])

		imgs[2] = k.Filters.Rotate(imgs[1])
		grays[2] = k.Filters.Gray(imgs[2])
		faces[2] = k.FaceDetector.Faces(grays[2])

		imgs[3] = k.Filters.Rotate(imgs[2])
		grays[3] = k.Filters.Gray(imgs[3])
		faces[3] = k.FaceDetector.Faces(grays[3])

		i = maxIndex(faces)
		k.Monitor.MsgWithDuration(time.Since(t), "%s detected faces: %s %s %s %s",
			photo.Name,
			picked(faces[0], i, 0),
			picked(faces[1], i, 1),
			picked(faces[2], i, 2),
			picked(faces[3], i, 3),
		)

		k.PhotoPersister.Persist(imgs[i], photo.Name)
		if err := k.Uploader.Upload(photo.Name); err != nil {
			time.Sleep(50 * time.Millisecond)
			k.Monitor.Msg("", "retrying to upload")
			k.Uploader.Upload(photo.Name)
		}
	}
}

func picked(faces, picked, i int) string {
	if picked == i {
		return color.HiGreenString("%d", faces)
	}

	return color.WhiteString("%d", faces)
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
