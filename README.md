# kmscan

Canon Lide dedicated wrapper for scanimage(1)
with image processing enhancements for scanning photos.

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

Run without installing in the system (pass arguments as `ARGS`):

    make run ARGS='-h'

Run tests:

    make run test
