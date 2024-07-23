# kmscan

Canon Lide dedicated wrapper for scanimage(1) with image processing enhancements.

## How to run

Pull the repository:

    $ git clone https://github.io/maicher/kmscan
    $ cd kmscan

Install and run:

    # make build install
    $ kmscan

More info:

    $ kmscan -h

Uninstall:

    # make uninstall

## Development

Pass arguments as `ARGS`:

    make run ARGS='-h'

Run tests:

    make run test

## Troubleshooting

If scanner stops responding, kill the scanning process:

    $ pkill scanimage

or restart the scanner:

    # usb Canon restart

You can find the implementation of the restart script here:
https://github.com/maicher/dotfiles/blob/master/lib/.local/bin/usb
