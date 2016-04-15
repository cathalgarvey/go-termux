# go-termux
by Cathal Garvey, Copyright 2016, released under GNU AGPLv3 or later

[![GoDoc][godoc-badge]][godoc]

## Termux-API package implemented in pure Go
**Warning: Most of these APIs are untested. Some, I can't even test as I use SMSSecure (no SmsInbox call data), or have no TTS engines. Others, I simply haven't had time. In particular, calls that pass input buffers are as yet still untested. Help/bugs/PRs welcome!**

Termux is a great terminal app for Android. It has several optional extension
apps, including Termux:API, which exposes a unix socket RPC interface for
making Android API calls to device hardware and UI details.

To use Termux:API's RPC you ordinarily need to install the termux-api tools in
termux and write shell scripts. I wanted, instead, to be able to use Go to write
typechecked, multithreaded pseudo-apps. Using the additional extension Termux:Widget,
these could then be even given desktop icons and used like other apps, albeit
with a dependency on Termux, Termux:API and Termux:Widget.

go-termux reimplements the C tool and wrapper scripts from `termux-api` so that,
in theory, you won't even need to install `termux-api` to use the Termux:API app.
In fact, I'm not even certain you'll need Termux to take advantage of Termux:API,
though Termux appears to have permissions to open Unix ports that aren't available
to the normal adb shell.

[godoc]: https://godoc.org/github.com/cathalgarvey/go-termux "GoDoc"
[godoc-badge]: https://godoc.org/github.com/cathalgarvey/go-termux?status.svg "GoDoc Badge"
