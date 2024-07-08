package main

import (
	_ "embed"
	"fmt"

	"github.com/maicher/kmscan/internal/options"
)

//go:embed internal/config/kmscanrc.example.toml
var kmscanrcExample string

func main() {
	opts := options.Parse()
	fmt.Printf("%+v\n", opts)
}
