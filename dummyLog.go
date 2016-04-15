// +build !debug

package termux

type ctx map[string]interface{}

// Debug is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func logdebug(m string, p ...ctx) {
}

// Info is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func loginfo(m string, p ...ctx) {
}

// Error is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func logerror(m string, p ...ctx) {
}
