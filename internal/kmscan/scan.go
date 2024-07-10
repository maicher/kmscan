package kmscan

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Scan struct {
	Image   image.Image
	Maximum image.Image
	Gray    *image.Gray
	Binary  *image.Gray
}

func NewScan(imagePath string) (*Scan, error) {
	s := Scan{}

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image: %s", err)
	}
	defer file.Close()

	if s.Image, err = png.Decode(file); err != nil {
		return nil, fmt.Errorf("error decoding image: %s", err)
	}

	return &s, nil
}
