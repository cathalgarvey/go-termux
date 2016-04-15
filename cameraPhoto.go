package termux

import (
	"bytes"
	"errors"
)

var (
	// ErrNoCameraSpecified is returned if cameraID is blank.
	ErrNoCameraSpecified = errors.New("Must specify a camera ID")
)

// CameraPhoto takes a photo from the specified camera and saves to outputFN.
// Size is undocumented in termux-api but appears to select a dimension value.
// You can leave Size empty.
func CameraPhoto(cameraID, outputFN, size string) error {
	var args []string
	if cameraID == "" {
		return ErrNoCameraSpecified
	}
	args = append(args, []string{"--es", "camera", cameraID}...)
	if outputFN == "" {
		return ErrNoOutputFile
	}
	args = append(args, []string{"--es", "file", "`realpath " + outputFN + "`"}...)
	if size != "" {
		args = append(args, []string{"--ei", "size_index", size}...)
	}
	_, err := toolExec(bytes.NewBuffer(nil), "CameraInfo", args...)
	return err
}
