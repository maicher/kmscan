package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var help = `NAME
  kmscan - Canon Lide dedicated wrapper for scanimage(1)
           with image processing enhancements for scanning photos

SYNOPSIS
  kmscan [OPTION...]

DESCRIPTION
  Once started, kmscan will wait for pressing the SEND button located on the
  scanner device.
  After clicking, it will perform scanning and auto extract the photos.

OPTIONS
  --config PATH,         -c   path to a kmscanrc.toml file
                              if not set, kmscan will to look up:
                              $XDG_CONFIG_HOME/kmscan/kmscanrc.toml,
                              $HOME/.config/kmscan/kmscanrc.toml,
                              or fallback to the defaults;
                              CLI options takes precedence over
  --detect,              -d   force detecting scanimage device
                              otherwise device will be read from cache
  --resolution NUM,      -r   resolution of the scan in dpi
                              available values: 300, 600 (default), 1200, 2400, 4800
  --min-width NUM,       -W   (default: 2000)
  --min-height NUM,      -H   (default: 2000)
  --min-aspect-ratio NUM -a   (default: 0.5)
  --max-aspect-ratio NUM -A   (default: 2.0)
  --brightness NUM,      -b   (default: 15)
  --window NUM,          -w   (default: 13)
  --threshold NUM,       -t   (default: 245)
  --debug,               -D   debug mode, which saves the filtered images into degug dir
  --debug-dir                 (default: debug) path to a dir to save the filtered images
  --image-path PATH,     -i   process the given image, instead of scanning it
  --cache-dir PATH            path to a cache dir
                              if not set, kmscan will try use:
                              $XDG_CACHE_HOME/kmscan or $HOME/.cache/kmscan
  --result-dir PATH           path to a result dir
                              (deault: ./results)
  --no-color,            -n   disable colored output
  --api-url,             -u   (optional) this option allows you to upload scanned images
                              to a server via an API endpoint,
                              when used, the program will send the image data
                              to the server in the body of a POST request,
                              the structure of the body is as follows:
                                filename: The name of the file being uploaded
                                data: The image file encoded as a Base64 string
  --api-key,             -k   (optional) API key to authenticate the upload request,
                              the key is sent as a HTTP header  under the "api-key" field
`

type Options struct {
	ConfigPath  string
	CacheDir    string
	ResultDir   string
	NoColor     bool
	ForceDetect bool
	Debug       bool
	DebugDir    string
	ImagePath   string
	APIURL      string
	APIKey      string

	Resolution     int
	MinHeight      int
	MinWidth       int
	MinAspectRatio float64
	MaxAspectRatio float64
	Brightness     float64
	Window         int
	Threshold      int
}

func ParseOptions() Options {
	var opts Options

	flag.StringVar(&opts.ConfigPath, "config", "", "")
	flag.StringVar(&opts.ConfigPath, "c", "", "")

	flag.StringVar(&opts.CacheDir, "cache-dir", "", "")
	flag.StringVar(&opts.CacheDir, "debug-dir", "", "")

	flag.StringVar(&opts.ResultDir, "result-dir", "results", "")

	flag.BoolVar(&opts.NoColor, "no-color", false, "")
	flag.BoolVar(&opts.NoColor, "n", false, "")

	flag.BoolVar(&opts.ForceDetect, "detect", false, "")
	flag.BoolVar(&opts.ForceDetect, "d", false, "")

	flag.BoolVar(&opts.Debug, "debug", false, "")
	flag.BoolVar(&opts.Debug, "D", false, "")

	flag.StringVar(&opts.ImagePath, "image-path", "", "")
	flag.StringVar(&opts.ImagePath, "i", "", "")

	flag.StringVar(&opts.APIURL, "api-url", "", "")
	flag.StringVar(&opts.APIURL, "u", "", "")

	flag.StringVar(&opts.APIKey, "api-key", "", "")
	flag.StringVar(&opts.APIKey, "k", "", "")

	flag.IntVar(&opts.Resolution, "resolution", 600, "")
	flag.IntVar(&opts.Resolution, "r", 600, "")

	flag.IntVar(&opts.MinWidth, "min-width", 2000, "")
	flag.IntVar(&opts.MinWidth, "W", 2000, "")

	flag.IntVar(&opts.MinHeight, "min-height", 2000, "")
	flag.IntVar(&opts.MinHeight, "H", 2000, "")

	flag.Float64Var(&opts.MinAspectRatio, "min-aspect-ratio", 0.5, "")
	flag.Float64Var(&opts.MinAspectRatio, "a", 0.5, "")

	flag.Float64Var(&opts.MaxAspectRatio, "max-aspect-ratio", 2.0, "")
	flag.Float64Var(&opts.MaxAspectRatio, "A", 2.0, "")

	flag.Float64Var(&opts.Brightness, "brightness", 2.0, "")
	flag.Float64Var(&opts.Brightness, "b", 2.0, "")

	flag.IntVar(&opts.Window, "window", 13, "")
	flag.IntVar(&opts.Window, "w", 13, "")

	flag.IntVar(&opts.Threshold, "threshold", 245, "")
	flag.IntVar(&opts.Threshold, "t", 245, "")

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprint(f, help) }
	flag.Parse()

	if opts.CacheDir == "" {
		if dir, ok := os.LookupEnv("XDG_CACHE_HOME"); ok {
			opts.CacheDir = filepath.Join(dir, "kmscan")
		} else if dir, ok := os.LookupEnv("HOME"); ok {
			opts.CacheDir = filepath.Join(dir, ".cache/kmscan")
		}
	}

	opts.DebugDir = "debug"

	return opts
}
