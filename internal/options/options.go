package options

import (
	"flag"
	"fmt"
)

var help = `NAME
  kmscan - Canon Lide dedicated wrapper for scanimage(1)
           with image processing enhancements  for scanning photos

SYNOPSIS
  kmscan [OPTION...]

DESCRIPTION
  Once started, kmscan will wait for pressing the SEND button located on the
  scanner device.
  After clicking, will perform scanning with predefined parametera
  and process the resulted image by extracting the photos and synchronizing
  them to the server.

OPTIONS
  --config PATH,    -c   path to a kmscanrc file
                         if not set, kmscan will try to look up following paths:
                         $XDG_CONFIG_HOME/kmscan/kmscanrc.toml
                         $HOME/.config/kmscan/kmscanrc.toml

CONFIG
  See the below link for example config:
    https://github.com/maicher/kmscan/blob/master/internal/config/kmscanrc.example.toml
`

type Options struct {
	ConfigPath string
}

func Parse() Options {
	var opts Options

	flag.StringVar(&opts.ConfigPath, "config", "", "")
	flag.StringVar(&opts.ConfigPath, "c", "", "")

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprint(f, help) }
	flag.Parse()

	return opts
}
