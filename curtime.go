package termux

import (
	"os/exec"
	"time"
)

// Curtime returns the current time according to the `date` Unix command.
// The time.Now() command seems to get it wrong sometimes, perhaps because
// of daylight savings, on Android.
func Curtime() (time.Time, error) {
	dateOutput, err := exec.Command("date").CombinedOutput()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.UnixDate, string(dateOutput))
}
