package termux

import (
	"encoding/json"
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
	return cameraPhoto(toolExec, cameraID, outputFN, size)
}

func cameraPhoto(execF toolExecFunc, cameraID, outputFN, size string) error {
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
	_, err := toolExec(nil, "CameraInfo", args...)
	return err
}

// CameraInfoResponse contains information on one camera.
type CameraInfoResponse struct {
	ID              string `json:"id"`
	Facing          string `json:"facing"`
	JpegOutputSizes []struct {
		Width  int
		Height int
	} `json:"jpeg_output_sizes"`
	FocalLengths      []float64 `json:"focal_lengths"`
	AutoExposureModes []string  `json:"auto_exposure_modes"`
	PhysicalSize      struct {
		Width  float64
		Height float64
	} `json:"physical_size"`
	Capabilities []string `json:"capabilities"`
}

// CameraInfo returns information on cameras in this device.
func CameraInfo() ([]CameraInfoResponse, error) {
	return cameraInfo(toolExec)
}

func cameraInfo(execF toolExecFunc) ([]CameraInfoResponse, error) {
	var resp []CameraInfoResponse
	ciBytes, err := execF(nil, "CameraInfo")
	if err != nil {
		return nil, err
	}
	if len(ciBytes) == 0 {
		return nil, ErrZeroLengthResponse
	}
	err = json.Unmarshal(ciBytes, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
