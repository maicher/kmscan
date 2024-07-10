package kmscan

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/maicher/kmscan/internal/monitor"
)

type Kmscan struct {
	Resolution     int
	Filters        *Filters
	Extractor      *Extractor
	DebugPersister Persister
	Monitor        *monitor.Monitor
}

func (k *Kmscan) Run(devices []string) {
	resultsCh := make(chan string, 10)

	for i, device := range devices {
		go k.loop(device, i+1, resultsCh)
	}

	for imagePath := range resultsCh {
		time.Sleep(10 * time.Millisecond)

		scan, err := NewScan(imagePath)
		if err != nil {
			k.Monitor.ProcessorError("%s\n", err)
			continue
		}
		defer os.Remove(imagePath)

		k.Filters.ApplyFilters(scan)
		if err := k.persist(scan); err != nil {
			k.Monitor.ProcessorError("%s\n", err)
		}

		k.Extractor.Extract(scan.Binary, scan.Image)
	}
}

func (k *Kmscan) ProcessImage(imagePath string) {
	scan, err := NewScan(imagePath)
	if err != nil {
		k.Monitor.ProcessorError("%s\n", err)
		return
	}

	k.Monitor.Scanner(imagePath, "file loaded")

	k.Filters.ApplyFilters(scan)
	if err := k.persist(scan); err != nil {
		k.Monitor.ProcessorError("%s\n", err)
	}

	k.Extractor.Extract(scan.Binary, scan.Image)
}

func (k *Kmscan) loop(device string, deviceNo int, resultsCh chan<- string) {
	var i int
	var calibrate string

	for {
		// Open a temporary file in RAM (in the /dev/shm directory)
		file, err := os.CreateTemp("/dev/shm", "scan_*.png")
		if err != nil {
			k.Monitor.ProcessorError("Failed to create temporary file: %s\n", err)
		}

		if i%10 == 0 {
			k.Monitor.Scanner("press SEND button on the scanner", "ready for scan %d with calibrate", i+1)
			calibrate = "always"
		} else {
			k.Monitor.Scanner("press SEND button on the scanner", "ready for scan %d", i+1)
			calibrate = "never"
		}

		cmd := exec.Command("scanimage",
			"-d", device,
			"--resolution", fmt.Sprintf("%ddpi", k.Resolution),
			"--mode", "Color",
			"--format", "png",
			"--calibrate", calibrate,
			"--button-controlled=yes",
		)

		cmd.Stdout = file
		if err = cmd.Run(); err != nil {
			k.Monitor.ProcessorError("error scanning: %s\n", err)
			return
		}
		k.Monitor.Scanner("", "done %d", i+1)

		resultsCh <- file.Name()
		i++

		time.Sleep(100 * time.Millisecond)
	}
}

func (k *Kmscan) persist(scan *Scan) error {
	if err := k.DebugPersister.Persist(scan.Image, "_0_scan.jpg"); err != nil {
		return err
	}

	if err := k.DebugPersister.Persist(scan.Maximum, "_1_max.jpg"); err != nil {
		return err
	}

	if err := k.DebugPersister.Persist(scan.Gray, "_2_gray.jpg"); err != nil {
		return err
	}

	if err := k.DebugPersister.Persist(scan.Binary, "_3_binary.jpg"); err != nil {
		return err
	}

	return nil
}
