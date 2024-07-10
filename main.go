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

const cmdName = "scanimage"

func main() {
	var err error

	opts := config.ParseOptions()
	m := monitor.New()

	resultPersister, err := kmscan.NewFilePersister(opts.ResultDir, m)
	if err != nil {
		m.InitError("%s", err)
		os.Exit(1)
	}

	k := kmscan.Kmscan{
		Resolution: opts.Resolution,
		Filters: &kmscan.Filters{
			Brightness: float32(opts.Brightness),
			Window:     opts.Window,
			Threshold:  opts.Threshold,
			Monitor:    m,
		},
		Extractor: &kmscan.Extractor{
			MinHeight:      opts.MinHeight,
			MinWidth:       opts.MinWidth,
			MinAspectRatio: opts.MinAspectRatio,
			MaxAspectRatio: opts.MaxAspectRatio,
			Persister:      resultPersister,
			Monitor:        m,
		},
		Monitor: m,
	}

	if opts.Debug {
		if k.DebugPersister, err = kmscan.NewFilePersister(opts.DebugDir, m); err != nil {
			m.InitError("%s", err)
			os.Exit(1)
		}
	} else {
		k.DebugPersister = kmscan.NullPersister{}
	}

	if opts.ImagePath != "" {
		k.ProcessImage(opts.ImagePath)

		return
	}

	d := detector.New(m, cmdName)
	devices, err := d.ReadOrDetect(opts.CacheDir, opts.ForceDetect)
	if err != nil {
		m.InitError("%s", err)
		os.Exit(1)
	}

	k.Run(devices)
}
