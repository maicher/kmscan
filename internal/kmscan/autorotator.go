package kmscan

import (
	"image"
	"time"

	"github.com/fatih/color"
	"github.com/maicher/kmscan/internal/ui"
)

type Autorotator struct {
	FaceDetector *FaceDetector
	Filters      *Filters
	Logger       *ui.Logger
}

func (a *Autorotator) Autorotate(img image.Image, name string) image.Image {
	t := time.Now()

	imgs := [4]image.Image{}
	grays := [4]*image.Gray{}
	faces := [4]int{}
	var picked int

	// Rotate the image 3 times and check in which orientation the most faces are found.
	imgs[0] = img
	for i := 0; i < 4; i++ {
		if i != 0 {
			imgs[i] = a.Filters.Rotate(imgs[i-1])
		}
		grays[i] = a.Filters.Gray(imgs[i])
		faces[i] = a.FaceDetector.Faces(grays[i])
	}

	picked = a.maxIndex(faces)
	a.Logger.MsgWithDuration(time.Since(t), "%s faces: %s %s %s %s",
		name,
		a.msg(faces[0], picked, 0),
		a.msg(faces[1], picked, 1),
		a.msg(faces[2], picked, 2),
		a.msg(faces[3], picked, 3),
	)

	return imgs[picked]
}
func (a *Autorotator) msg(faces, picked, i int) string {
	if picked == i {
		return color.HiWhiteString("%d", faces)
	}

	return color.WhiteString("%d", faces)
}

func (a *Autorotator) maxIndex(arr [4]int) int {
	maxIdx := 0
	maxValue := arr[0]

	for i, value := range arr {
		if value > maxValue {
			maxValue = value
			maxIdx = i
		}
	}

	return maxIdx
}
