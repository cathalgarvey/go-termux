package termux

import "strconv"

// Vibrate vibrates the phone for the desired milliseconds
func Vibrate(msDuration int, force bool) error {
	var args []string
	if msDuration == 0 {
		return nil
	}
	args = append(args, []string{"--ei", "duration_ms", strconv.Itoa(msDuration)}...)
	if force {
		args = append(args, []string{"--ez", "force", "true"}...)
	}
	_, err := toolExec(nil, "Vibrate", args...)
	return err
}
