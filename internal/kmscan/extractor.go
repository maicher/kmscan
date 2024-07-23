package kmscan

import (
	"fmt"
	"image"
	"time"

	"github.com/maicher/kmscan/internal/monitor"
)

type Extractor struct {
	MinWidth       int
	MinHeight      int
	MinAspectRatio float64
	MaxAspectRatio float64

	Monitor   *monitor.Monitor
	Persister *FilePersister
}

// Extracts individual photos from the scanning result.
func (e *Extractor) Extract(scan *Scan, out chan<- Photo) {
	// Detect the rectangles
	rectangles := e.findRectangles(scan.Binary)

	// Filter rectangles based on size and aspect ratio criteria
	filteredRectangles := e.filterRectangles(rectangles)

	// Save each detected rectangle as an image
	t := time.Now().Format("20060102_150405")
	for i, rect := range filteredRectangles {
		roi := scan.Image.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)

		out <- Photo{
			Name:  fmt.Sprintf("%s_%02d.jpg", t, i),
			Image: roi,
		}
	}
}

func (e *Extractor) findRectangles(img *image.Gray) []image.Rectangle {
	t := time.Now()

	bounds := img.Bounds()
	var rectangles []image.Rectangle
	visited := make([][]bool, bounds.Max.Y)
	for i := range visited {
		visited[i] = make([]bool, bounds.Max.X)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y == 0 && !visited[y][x] {
				rect := e.findBoundingBox(img, x, y, visited)
				rectangles = append(rectangles, rect)
			}
		}
	}

	e.Monitor.MsgWithDuration(time.Since(t), "%d rectangles found", len(rectangles))

	return rectangles
}

func (e *Extractor) findBoundingBox(img *image.Gray, startX, startY int, visited [][]bool) image.Rectangle {
	minX, minY := startX, startY
	maxX, maxY := startX, startY

	queue := []image.Point{{X: startX, Y: startY}}
	visited[startY][startX] = true

	dirs := []image.Point{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		for _, dir := range dirs {
			np := image.Point{X: p.X + dir.X, Y: p.Y + dir.Y}

			if np.X >= 0 && np.X < img.Bounds().Max.X && np.Y >= 0 && np.Y < img.Bounds().Max.Y &&
				img.GrayAt(np.X, np.Y).Y == 0 && !visited[np.Y][np.X] {
				queue = append(queue, np)
				visited[np.Y][np.X] = true

				if np.X < minX {
					minX = np.X
				}
				if np.X > maxX {
					maxX = np.X
				}
				if np.Y < minY {
					minY = np.Y
				}
				if np.Y > maxY {
					maxY = np.Y
				}
			}
		}
	}

	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func (e *Extractor) filterRectangles(rectangles []image.Rectangle) []image.Rectangle {
	t := time.Now()

	var filtered []image.Rectangle
	for i, rect := range rectangles {
		width := rect.Dx()
		height := rect.Dy()
		aspectRatio := float64(width) / float64(height)

		if width >= e.MinWidth || height >= e.MinHeight {
			if aspectRatio >= e.MinAspectRatio && aspectRatio <= e.MaxAspectRatio {
				e.Monitor.Msg(fmt.Sprintf("ratio %0.2f", aspectRatio), "rectangle %03d picked", i)

				filtered = append(filtered, rect)
			} else {
				e.Monitor.Msg(fmt.Sprintf("ration %0.2f", aspectRatio), "rectangle %03d rejected", i)
			}

		}
	}

	e.Monitor.MsgWithDuration(time.Since(t), "%d picked", len(filtered))

	return filtered
}
