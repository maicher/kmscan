package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/guregu/null.v3"
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

CONFIGURATION
  kmscan can be configured by kmscanrc file.
  See default: https://github.com/maicher/kmscan/blob/master/internal/config/kmscanrc.example.toml
  Copy the above file into:
    $XDG_CONFIG_HOME/kmscan/kmscanrc.toml and modify.
  or, of $XDG_CONFIG_HOME is not defined, into:
    $HOME/.config/kmscan/kmscanrc.toml
  CLI options has precedence over options defined in the kmscanrc file.

OPTIONS
  Main options (all options are optional):

  --profile-name NAME,   -p   scanning profile name, when not specified the first will be used
                              profiles are defined in the example configuration file
                              (see CONFIGURATION section for more info)
  --cache-dir PATH            path to a cache dir
                              (default: $XDG_CACHE_HOME/kmscan or $HOME/.cache/kmscan)
  --result-dir PATH           path to a result dir
                              (deault: ./results)
  --force-detect,        -d   force detecting scanimage device;
                              kmscan, when run for the first time, will detect scanimage device and store
                              this information in the cache dir;
                              set this flag if you want to ignore the cache information and re-detect the device again
  --debug,               -D   debug mode
                              kmscan after performing the scanning, applies filters to be able to extract individual
                              photos from the scanned image; with this options the scanned image with
                              applied filters will be saved to the debug dir so that it can be visually
                              checked by the user
  --debug-dir                 path to a dir to save the images with applied filters
                              (default: ./debug) this option is ignored when -D flag is not provided
  --no-color,            -n   disable colored output
  --image-path PATH,     -i   process the given image, instead of scanning image (useful for debugging)
  --api-url,             -u   (optional) this option allows you to upload scanned images
                              to a server via an API endpoint;
                              when used, kmscan will send the image
                              to the server in the body of a POST request;
                              the structure of the body:
                                filename: The name of the file being uploaded
                                data: The image file encoded as a Base64 string
  --api-key,             -k   (optional) API key to authenticate the upload request,
                              the key is sent as a HTTP header under the "api-key" field
`

type Options struct {
	ProfileName string
	CacheDir    string
	ResultDir   string
	ForceDetect bool
	Debug       bool
	DebugDir    string
	NoColor     bool
	ImagePath   string
	APIURL      string
	APIKey      string
}

func ParseOptions(optsRC *OptionsRC) Options {
	var opts Options

	// Main Options.
	flag.StringVar(&opts.ProfileName, "profile-name", optsRC.ProfileName.String, "")
	flag.StringVar(&opts.ProfileName, "p", optsRC.ProfileName.String, "")

	flag.StringVar(&opts.CacheDir, "cache-dir", optsRC.CacheDir.String, "")

	flag.StringVar(&opts.ResultDir, "result-dir", withDefault(optsRC.ResultDir, "results"), "")

	flag.BoolVar(&opts.ForceDetect, "force-detect", withDefaultBool(optsRC.ForceDetect, false), "")
	flag.BoolVar(&opts.ForceDetect, "d", withDefaultBool(optsRC.ForceDetect, false), "")

	flag.BoolVar(&opts.Debug, "debug", withDefaultBool(optsRC.Debug, false), "")
	flag.BoolVar(&opts.Debug, "D", withDefaultBool(optsRC.Debug, false), "")

	flag.StringVar(&opts.DebugDir, "debug-dir", withDefault(optsRC.DebugDir, "debug"), "")

	flag.BoolVar(&opts.NoColor, "no-color", false, "")
	flag.BoolVar(&opts.NoColor, "n", false, "")

	flag.StringVar(&opts.ImagePath, "image-path", optsRC.ImagePath.String, "")
	flag.StringVar(&opts.ImagePath, "i", optsRC.ImagePath.String, "")

	flag.StringVar(&opts.APIURL, "api-url", optsRC.APIURL.String, "")
	flag.StringVar(&opts.APIURL, "u", optsRC.APIURL.String, "")

	flag.StringVar(&opts.APIKey, "api-key", optsRC.APIKey.String, "")
	flag.StringVar(&opts.APIKey, "k", optsRC.APIKey.String, "")

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

	return opts
}

func withDefault(val null.String, defaultVal string) string {
	if val.Valid {
		return val.String
	}

	return defaultVal
}

func withDefaultBool(val null.Bool, defaultVal bool) bool {
	if val.Valid {
		return val.Bool
	}

	return defaultVal
}
