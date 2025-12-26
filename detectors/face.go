package detectors

import (
	"fmt"
	"image"
	"log"
	"time"

	"gocv.io/x/gocv"

	"FacialRecognition/consts"
	"FacialRecognition/logging"
	"FacialRecognition/misc"
)

// facialDetector holds all classifiers and detection parameters
type facialDetector struct {
	frontalFaceClassifier *gocv.CascadeClassifier
	sideFaceClassifier    *gocv.CascadeClassifier
	eyeClassifier         *gocv.CascadeClassifier

	// Detection parameters
	faceParams detectionParams
	eyeParams  detectionParams

	// Tracking how long a face has been there
	lastFaceDetectedTime time.Time
	faceAbsentLogged     bool
	faceTimeout          time.Duration
}

type detectionParams struct {
	ScaleFactor  float64
	MinNeighbors int
	Flags        int
	MinSize      image.Point
	MaxSize      image.Point
}

// NewFacialDetector creates and initializes all classifiers once
func newFacialDetector() (*facialDetector, error) {
	fd := &facialDetector{
		// Default face detection parameters
		faceParams: detectionParams{
			ScaleFactor:  1.2,
			MinNeighbors: 4,
			Flags:        0,
			MinSize:      image.Point{150, 150},
			MaxSize:      image.Point{},
		},
		// Default eye detection parameters
		eyeParams: detectionParams{
			ScaleFactor:  1.2,
			MinNeighbors: 7,
			Flags:        0,
			MinSize:      image.Point{25, 15},
			MaxSize:      image.Point{90, 70},
		},
		// Initialize time Tracking
		faceTimeout: 30 * time.Second,
	}

	// Load frontal face classifier
	fd.frontalFaceClassifier = &gocv.CascadeClassifier{}
	*fd.frontalFaceClassifier = gocv.NewCascadeClassifier()
	if !fd.frontalFaceClassifier.Load(consts.PATH_TO_FRONT_FACE_CASCADE) {
		return nil, fmt.Errorf("error loading frontal face cascade: %s", consts.PATH_TO_FRONT_FACE_CASCADE)
	}

	// Load side face classifier
	fd.sideFaceClassifier = &gocv.CascadeClassifier{}
	*fd.sideFaceClassifier = gocv.NewCascadeClassifier()
	if !fd.sideFaceClassifier.Load(consts.PATH_TO_SIDE_FACE_CASCADE) {
		fd.frontalFaceClassifier.Close()
		return nil, fmt.Errorf("error loading side face cascade: %s", consts.PATH_TO_SIDE_FACE_CASCADE)
	}

	// Load eye classifier
	fd.eyeClassifier = &gocv.CascadeClassifier{}
	*fd.eyeClassifier = gocv.NewCascadeClassifier()
	if !fd.eyeClassifier.Load(consts.PATH_TO_EYE_CASCADE) {
		fd.frontalFaceClassifier.Close()
		fd.sideFaceClassifier.Close()
		return nil, fmt.Errorf("error loading eye cascade: %s", consts.PATH_TO_EYE_CASCADE)
	}

	return fd, nil
}

// Close releases all classifiers
func (fd *facialDetector) Close() {
	if fd.frontalFaceClassifier != nil {
		fd.frontalFaceClassifier.Close()
	}
	if fd.sideFaceClassifier != nil {
		fd.sideFaceClassifier.Close()
	}
	if fd.eyeClassifier != nil {
		fd.eyeClassifier.Close()
	}
}

// DetectionResult holds all detected features
type detectionResult struct {
	FrontalFaces []image.Rectangle
	SideFaces    []image.Rectangle
	Eyes         []image.Rectangle
	ValidEyes    []image.Rectangle // Eyes inside faces
}

// DetectAll performs all detections in one pass
func (fd *facialDetector) detectAll(img gocv.Mat) detectionResult {
	result := detectionResult{}

	// Detect frontal faces
	result.FrontalFaces = fd.frontalFaceClassifier.DetectMultiScaleWithParams(
		img,
		fd.faceParams.ScaleFactor,
		fd.faceParams.MinNeighbors,
		fd.faceParams.Flags,
		fd.faceParams.MinSize,
		fd.faceParams.MaxSize,
	)

	// Detect side faces
	result.SideFaces = fd.sideFaceClassifier.DetectMultiScaleWithParams(
		img,
		fd.faceParams.ScaleFactor,
		fd.faceParams.MinNeighbors,
		fd.faceParams.Flags,
		fd.faceParams.MinSize,
		fd.faceParams.MaxSize,
	)

	// Detect eyes
	result.Eyes = fd.eyeClassifier.DetectMultiScaleWithParams(
		img,
		fd.eyeParams.ScaleFactor,
		fd.eyeParams.MinNeighbors,
		fd.eyeParams.Flags,
		fd.eyeParams.MinSize,
		fd.eyeParams.MaxSize,
	)

	// Find valid eyes (inside faces)
	result.ValidEyes = fd.filterEyesInFaces(result.Eyes, result.FrontalFaces, result.SideFaces)

	return result
}

// filterEyesInFaces returns only eyes that are inside detected faces
func (fd *facialDetector) filterEyesInFaces(eyes, frontalFaces, sideFaces []image.Rectangle) []image.Rectangle {
	var validEyes []image.Rectangle

	for _, eye := range eyes {
		// Check if eye is in any frontal face
		for _, face := range frontalFaces {
			if misc.RectIsContained(eye, face) {
				validEyes = append(validEyes, eye)
				goto nextEye // Skip to next eye once we found it's valid
			}
		}

		// Check if eye is in any side face
		for _, face := range sideFaces {
			if misc.RectIsContained(eye, face) {
				validEyes = append(validEyes, eye)
				break
			}
		}
	nextEye:
	}

	return validEyes
}

// DrawDetections draws all detected features on the image
func (fd *facialDetector) DrawDetections(img *gocv.Mat, result detectionResult) {
	// Draw frontal faces
	for _, face := range result.FrontalFaces {
		gocv.Rectangle(img, face, consts.FACE_COLOR, 3)
	}

	// Draw side faces
	for _, face := range result.SideFaces {
		gocv.Rectangle(img, face, consts.FACE_COLOR, 3)
	}

	// Draw valid eyes only
	for _, eye := range result.ValidEyes {
		gocv.Rectangle(img, eye, consts.EYE_COLOR, 3)
	}
}

// ProcessFrame handles a single frame
func (fd *facialDetector) processFrame(img *gocv.Mat) {
	// for some strange reason this line is required for the code to detect the logging package
	// fmt.Printf("Type of logging package: %T\n", logging.LogDetectedFace)
	result := fd.detectAll(*img)

	// Track time for faces
	// Check if any faces were detected
	facesDetected := len(result.FrontalFaces) > 0 || len(result.SideFaces) > 0
	currentTime := time.Now()

	if facesDetected {
		// Face detected
		if fd.faceAbsentLogged {
			// Face just came back - log it here with your existing logging
			// Call your logging function here
			log.Printf("Face reappeared after absence, logging")

			err := logging.LogDetectedFace()
			if err != nil {
				log.Printf("Could not access logging file, aborting log operation: %v", err)
			} else {
				log.Println("Writing to file successful")
			}
			fd.faceAbsentLogged = false
		}
		fd.lastFaceDetectedTime = currentTime
	} else {
		// No face detected
		if !fd.lastFaceDetectedTime.IsZero() {
			timeSinceFace := currentTime.Sub(fd.lastFaceDetectedTime)
			if timeSinceFace >= fd.faceTimeout && !fd.faceAbsentLogged {
				log.Printf("Face not detected for %v", fd.faceTimeout)
				fd.faceAbsentLogged = true
			}
		}
	}

	fmt.Printf("Found %d front face(s), %d side face(s), %d eye(s) (%d valid)\n",
		len(result.FrontalFaces),
		len(result.SideFaces),
		len(result.Eyes),
		len(result.ValidEyes))

	fd.DrawDetections(img, result)
}

// Main detection loop
func DetectFacialFeatures(webcam *gocv.VideoCapture) {
	// Initialize detector once
	detector, err := newFacialDetector()
	if err != nil {
		log.Fatal(err)
	}
	defer detector.Close()

	window := gocv.NewWindow("Face Detection")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Println("Cannot read from webcam")
			break
		}

		if img.Empty() {
			continue
		}

		// Process the frame
		detector.processFrame(&img)

		// Show the image
		window.IMShow(img)

		// Exit on ESC
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
