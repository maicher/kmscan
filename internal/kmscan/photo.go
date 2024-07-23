package kmscan

import "image"

type Photo struct {
	Name  string
	Image image.Image
	Gray  *image.Gray
}
