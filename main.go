package main

import (
	_ "embed"
	"os"

	"github.com/fatih/color"
	"github.com/maicher/kmscan/internal/config"
	"github.com/maicher/kmscan/internal/detector"
	"github.com/maicher/kmscan/internal/kmscan"
	"github.com/maicher/kmscan/internal/ui"
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
	if opts.NoColor {
		color.NoColor = true
	}

	mainLogger := ui.NewMainLogger()
	procLogger := ui.NewProcLogger()
	scanLogger := ui.NewScanLogger()
	apiLogger := ui.NewAPILogger()
	uiLogger := ui.NewUILogger()

	filters := &kmscan.Filters{
		Brightness: float32(opts.Brightness),
		Window:     opts.Window,
		Threshold:  opts.Threshold,
		Logger:     procLogger,
	}

	faceDetector := &kmscan.FaceDetector{
		Puploc:     puploc,
		Facefinder: facefinder,
		Logger:     procLogger,
	}
	if err := faceDetector.Load(); err != nil {
		mainLogger.Err("unable to load face detector: %s", err)
		os.Exit(1)
	}

	k := kmscan.Kmscan{
		Scanner: &kmscan.Scanner{
			Resolution: opts.Resolution,
			Logger:     scanLogger,
		},
		Filters: filters,
		Extractor: &kmscan.Extractor{
			MinHeight:      opts.MinHeight,
			MinWidth:       opts.MinWidth,
			MinAspectRatio: opts.MinAspectRatio,
			MaxAspectRatio: opts.MaxAspectRatio,
			Logger:         procLogger,
		},
		Autorotator: &kmscan.Autorotator{
			FaceDetector: faceDetector,
			Filters:      filters,
			Logger:       procLogger,
		},
		Uploader: kmscan.NewUploader(opts.ResultDir, opts.APIURL, opts.APIKey, apiLogger),
		Keyboard: &ui.Keyboard{
			Logger: uiLogger,
		},
		Logger: mainLogger,
	}

	if k.PhotoPersister, err = kmscan.NewFilePersister(opts.ResultDir, procLogger); err != nil {
		mainLogger.Err("%s", err)
		os.Exit(1)
	}

	if opts.Debug {
		if k.ScanPersister, err = kmscan.NewFilePersister(opts.DebugDir, procLogger); err != nil {
			mainLogger.Err("%s", err)
			os.Exit(1)
		}
	} else {
		k.ScanPersister = kmscan.NullPersister{}
	}

	if opts.ImagePath != "" {
		k.ProcessImage(opts.ImagePath)

		return
	}

	d := detector.New(mainLogger, cmdName)
	devices, err := d.ReadOrDetect(opts.CacheDir, opts.ForceDetect)
	if err != nil {
		mainLogger.Err("%s", err)
		os.Exit(1)
	}

	k.Run(devices)
}
