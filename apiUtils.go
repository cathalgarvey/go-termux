package termux

import (
	"net"
	"os"
	"strings"

	"github.com/twinj/uuid"
)

// Helper func that gets Env but omits any PATH entries.
func pathlessEnv() []string {
	var output []string
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "PATH=") {
			continue
		}
		output = append(output, entry)
	}
	return output
}

func getUnixSocket() net.UnixAddr {
	return net.UnixAddr{Name: "@" + uuid.NewV4().String(), Net: "unix"}
}

// Create a callback that checks errors and returns bool, but also appends
// errors to a slice if they are non-nil.
func makeErrCallback(funcName string, errSlice []error) func(error) bool {
	return func(err error) bool {
		if err != nil {
			logerror("Error in socket thread", ctx{"function": funcName, "error": err.Error()})
			errSlice = append(errSlice, err)
			return true
		}
		return false
	}
}
