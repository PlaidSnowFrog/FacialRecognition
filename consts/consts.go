package consts

import (
	"image/color"
)

// File paths and directories
const PATH_TO_FRONT_FACE_CASCADE = "data/haarcascade_frontalface_default.xml"
const PATH_TO_SIDE_FACE_CASCADE = "data/haarcascade_profileface.xml"
const PATH_TO_EYE_CASCADE = "data/haarcascade_eye.xml"
const PATH_TO_LOG_FILE = "log.txt"

// Colors
var FACE_COLOR = color.RGBA{0, 255, 0, 50}
var EYE_COLOR = color.RGBA{255, 0, 0, 50}
