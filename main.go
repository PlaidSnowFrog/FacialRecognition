package main

import (
	"fmt"

	"FacialRecognition/detectors"

	"gocv.io/x/gocv"
)

func main() {
	// set to use a video capture device 0
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// Detection
	detectors.DetectFacialFeatures(webcam)

	// show the image in the window, and wait 1 millisecond
	window.IMShow(img)
	window.WaitKey(1)
}
