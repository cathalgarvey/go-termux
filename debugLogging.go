// +build debug

package termux

import "gopkg.in/inconshreveable/log15.v2"

type ctx map[string]interface{}

// Debug is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func logdebug(m string, p ...ctx) {
	if len(p) > 0 {
		log15.Debug(m, log15.Ctx(p[0]))
	} else {
		log15.Debug(m)
	}
}

// Info is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func loginfo(m string, p ...ctx) {
	if len(p) > 0 {
		log15.Info(m, log15.Ctx(p[0]))
	} else {
		log15.Info(m)
	}
}

// Error is a shim for Log15, so that a build tag can be used to enable or
// disable logging at build-time.
func logerror(m string, p ...ctx) {
	if len(p) > 0 {
		log15.Error(m, log15.Ctx(p[0]))
	} else {
		log15.Error(m)
	}
}
