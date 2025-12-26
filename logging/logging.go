package logging

import (
	"fmt"
	"os"
	"time"

	"FacialRecognition/consts"
)

func LogDetectedFace() error {
	logFile, err := os.OpenFile(consts.PATH_TO_LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error loading not open log file: %s", consts.PATH_TO_LOG_FILE)
	}
	defer logFile.Close() // Add this to close the file
	currentTime := time.Now()
	// This line is hella long but it should work and its kinda simple
	// so don't fix it pls
	logFile.Write([]byte(fmt.Sprint("Face Detected in Doorway; Time: ", currentTime.Format(time.Stamp), "\n")))
	return nil
}
