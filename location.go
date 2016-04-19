package termux

import "encoding/json"

// LocationResponse is the JSON returned by the termux-location tool
type LocationResponse struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	Accuracy  float64
	Bearing   float64
	Speed     float64
	ElapsedMs int
	Provider  string
}

// Location makes a request for location information
func Location(request, provider string) (*LocationResponse, error) {
	return location(toolExec, request, provider)
}

func location(execF toolExecFunc, request, provider string) (*LocationResponse, error) {
	err := validateLocationArgs(request, provider)
	if err != nil {
		return nil, err
	}
	// ---
	var (
		args []string
		resp LocationResponse
	)
	args = append(args, []string{"--es", "request", request}...)
	args = append(args, []string{"--es", "provider", provider}...)
	locBytes, err := execF(nil, "Location", args...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(locBytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func validateLocationArgs(request, provider string) error {
	switch provider {
	case "gps", "network", "passive":
	default:
		return ErrBadLocationProvider
	}
	switch request {
	case "once", "last":
	case "updates":
		return ErrUpdatesNotSupportedYet
	default:
		return ErrBadLocationRequest
	}
	return nil
}
