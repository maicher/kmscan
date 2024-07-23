package main

import (
	_ "embed"
	"os"

	"github.com/maicher/kmscan/internal/config"
	"github.com/maicher/kmscan/internal/detector"
	"github.com/maicher/kmscan/internal/kmscan"
	"github.com/maicher/kmscan/internal/monitor"
)

//go:embed internal/config/kmscanrc.example.toml
var kmscanrcExample string

// Copied from: https://github.com/esimov/pigo/tree/master/cascade
//
//go:embed internal/config/puploc
var puploc []byte

// Copied from: https://github.com/esimov/pigo/tree/master/cascade
//
//go:embed internal/config/facefinder
var facefinder []byte

const cmdName = "scanimage"

func main() {
	var err error

	opts := config.ParseOptions()
	mon := monitor.NewMain()
	procMon := monitor.NewProcessor()

	k := kmscan.Kmscan{
		Scanner: &kmscan.Scanner{
			Resolution: opts.Resolution,
			Monitor:    monitor.NewScanner(),
		},
		Filters: &kmscan.Filters{
			Brightness: float32(opts.Brightness),
			Window:     opts.Window,
			Threshold:  opts.Threshold,
			Monitor:    procMon,
		},
		Extractor: &kmscan.Extractor{
			MinHeight:      opts.MinHeight,
			MinWidth:       opts.MinWidth,
			MinAspectRatio: opts.MinAspectRatio,
			MaxAspectRatio: opts.MaxAspectRatio,
			Monitor:        procMon,
		},
		FaceDetector: &kmscan.FaceDetector{
			Puploc:     puploc,
			Facefinder: facefinder,
			Monitor:    procMon,
		},
		Uploader:           kmscan.NewUploader(opts.ResultDir, monitor.NewAPI()),
		KeyboardController: &kmscan.KeyboardController{Monitor: monitor.NewUI()},
		Monitor:            mon,
	}

	if err := k.FaceDetector.Load(); err != nil {
		mon.Err("unable to load face detector: %s", err)
		os.Exit(1)
	}

	if k.PhotoPersister, err = kmscan.NewFilePersister(opts.ResultDir, procMon); err != nil {
		mon.Err("%s", err)
		os.Exit(1)
	}

	if opts.Debug {
		if k.ScanPersister, err = kmscan.NewFilePersister(opts.DebugDir, procMon); err != nil {
			mon.Err("%s", err)
			os.Exit(1)
		}
	} else {
		k.ScanPersister = kmscan.NullPersister{}
	}

	if opts.ImagePath != "" {
		k.ProcessImage(opts.ImagePath)

		return
	}

	d := detector.New(mon, cmdName)
	devices, err := d.ReadOrDetect(opts.CacheDir, opts.ForceDetect)
	if err != nil {
		mon.Err("%s", err)
		os.Exit(1)
	}

	k.Run(devices)
}
