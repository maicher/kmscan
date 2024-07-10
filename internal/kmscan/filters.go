package kmscan

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/disintegration/gift"
	"github.com/maicher/kmscan/internal/monitor"
)

type Filters struct {
	// Brightness correction percentage.
	Brightness float32

	// Window size for the Maximum rank order filter.
	// Should be an odd number, eg 9, 13, 15.
	Window int

	// Threshold for binary (black and white) conversion.
	// Should be between 0 and 255.
	Threshold int

	Monitor *monitor.Monitor
}

func (f *Filters) ApplyFilters(scan *Scan) {
	t := time.Now()

	scan.Maximum = f.max(scan.Image)
	scan.Gray = f.gray(scan.Maximum)
	scan.Binary = f.binary(scan.Gray)

	f.Monitor.Processor(time.Since(t), "filters applied")
}

func (f *Filters) gray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	return dst
}

func (f *Filters) max(src image.Image) image.Image {
	g := gift.New(
		gift.Brightness(15),
		gift.Maximum(13, true),
	)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	return dst
}

func (f *Filters) binary(src *image.Gray) *image.Gray {
	threshold := uint8(245)

	bounds := src.Bounds()
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if src.GrayAt(x, y).Y > threshold {
				dst.SetGray(x, y, color.Gray{Y: 255})
			} else {
				dst.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return dst
}
