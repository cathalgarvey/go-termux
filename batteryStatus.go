package termux

import (
	"bytes"
	"encoding/json"
)

// BatteryStatusResponse is returned by the BatteryStatus function.
type BatteryStatusResponse struct {
	Health      string
	Percentage  int
	Plugged     string
	Status      string
	Temperature float64
}

// BatteryStatus mimics the termux-battery-status script/call.
func BatteryStatus() (*BatteryStatusResponse, error) {
	var resp BatteryStatusResponse
	bsrBytes, err := toolExec(bytes.NewBuffer(nil), "BatteryStatus")
	if err != nil {
		return nil, err
	}
	if len(bsrBytes) == 0 {
		return nil, ErrZeroLengthResponse
	}
	err = json.Unmarshal(bsrBytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
