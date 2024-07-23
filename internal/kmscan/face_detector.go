package kmscan

import (
	"fmt"
	"image"

	pigo "github.com/esimov/pigo/core"
	"github.com/maicher/kmscan/internal/monitor"
)

type FaceDetector struct {
	Monitor    *monitor.Monitor
	Puploc     []byte
	Facefinder []byte

	face   *pigo.Pigo
	pupils *pigo.PuplocCascade
}

// Load face detection and pupil localization classifiers..
func (fd *FaceDetector) Load() error {
	var err error

	p := pigo.NewPigo()

	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	fd.face, err = p.Unpack(fd.Facefinder)
	if err != nil {
		return fmt.Errorf("Error unpacking the cascade file: %s", err)
	}

	fd.pupils, err = fd.pupils.UnpackCascade(fd.Puploc)
	if err != nil {
		return fmt.Errorf("Error unpacking the puploc cascade file: %s", err)
	}

	return nil
}

// Returns true if a face was found in the image.
func (fd *FaceDetector) Faces(img *image.Gray) int {
	// var face, left, right bool
	// t := time.Now()
	// var face bool
	var faces int

	imageParams := &pigo.ImageParams{
		Pixels: img.Pix,
		Rows:   img.Rect.Dy(),
		Cols:   img.Rect.Dx(),
		Dim:    img.Rect.Dx(),
	}

	results := fd.clusterDetection(imageParams)
	for i := 0; i < len(results); i++ {
		if results[i].Q < 20 {
			continue
		}

		faces++
		// faces++
		// // left eye
		// puploc := &pigo.Puploc{
		// 	Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
		// 	Col:      results[i].Col - int(0.185*float32(results[i].Scale)),
		// 	Scale:    float32(results[i].Scale) * 0.4,
		// 	Perturbs: 50,
		// }

		// det := fd.pupils.RunDetector(*puploc, *imageParams, 0, false)
		// // fmt.Printf("Left eye: %+v\n", det)
		// if det != nil {
		// 	left = true
		// }

		// // right eye
		// puploc = &pigo.Puploc{
		// 	Row:      results[i].Row - int(0.085*float32(results[i].Scale)),
		// 	Col:      results[i].Col + int(0.185*float32(results[i].Scale)),
		// 	Scale:    float32(results[i].Scale) * 0.4,
		// 	Perturbs: 50,
		// }
		// det = fd.pupils.RunDetector(*puploc, *imageParams, 0.0, false)
		// // fmt.Printf("Right eye: %+v\n", det)
		// if det != nil {
		// 	right = true
		// }
	}

	// return face && right && left

	// fd.Monitor.Processor(time.Since(t), "detected %d faces", faces)

	return faces
}

// clusterDetection runs Pigo face detector core methods
// and returns a cluster with the detected faces coordinates.
func (fd *FaceDetector) clusterDetection(imageParams *pigo.ImageParams) []pigo.Detection {
	cParams := pigo.CascadeParams{
		MinSize:     60,
		MaxSize:     600,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: *imageParams,
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.

	// Try with different angles.
	//angle := -1.571
	angle := 0.0
	dets := fd.face.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = fd.face.ClusterDetections(dets, angle)

	return dets
}
