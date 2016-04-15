package termux

import (
	"bytes"
	"encoding/json"
)

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
	var resp []CameraInfoResponse
	ciBytes, err := toolExec(bytes.NewBuffer(nil), "CameraInfo")
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
