package termux

import "encoding/json"

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
	return batteryStatus(toolExec)
}

func batteryStatus(execF toolExecFunc) (*BatteryStatusResponse, error) {
	var resp BatteryStatusResponse
	bsrBytes, err := execF(nil, "BatteryStatus")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bsrBytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
