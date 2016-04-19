/*Package termux provides a clean Go API for making calls to Termux's Android API tools.

It attempts to reimplement and replace the C code from termux-api, and to provide
a clean functional interface to all the api script wrappers, so that development
of functional termux "apps" can be as easy as:

		1. Install Termux
		2. Install Termux:API (just the app: no need to install termux-api!)
		3. Write and test your go app from Termux shell (possibly even compiling from Termux!)
		4. Optional: Use Termux:Widget to create Android desktop shortcuts for your new termux/Go "app"

At present, applications written using this framework appear to only function in
Termux shell, not in the ADB shell or system shell. It's possible this is because
Termux:API is unable to connect to unix ports outside of Termux's sandbox, I'm
not sure. Hopefully this can be fixed so that go-termux can be used for general
phone automation. Still more hopefully, perhaps Termux itself would prove unnecessary
and only Termux:API would be required in this case?

*/
package termux

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	am "github.com/cathalgarvey/androidam"
)

// toolExecFunc represents the interface all of the other tools construct args
// for and pass-through. It's interfaced here to allow for testing.
type toolExecFunc func(stdin io.Reader, tool string, toolargs ...string) ([]byte, error)

// Sends an intent broadcast for the target tool, providing arguments and
// input/output sockets for the API server.
// Rewrite of exec_am_broadcast from:
// https://github.com/termux/termux-packages/blob/master/packages/termux-api/termux-api.c#L27
func execAMBroadcast(inputAddress, outputAddress string, tool string, argv ...string) error {
	outputAddress = strings.Trim(outputAddress, "@")
	inputAddress = strings.Trim(inputAddress, "@")
	logdebug("In execAMBroadcast, prepping call", ctx{"inputAddress": inputAddress, "outputAddress": outputAddress, "argv": argv})
	bc := am.Broadcast(
		nil,
		&am.BroadcastArgs{
			User: am.GetCurrentAndroidUserID(),
		},
		&am.IntentArgs{
			Component: "com.termux.api/.TermuxApiReceiver",
			ExtraKeyStrings: map[string]string{
				"socket_input":  outputAddress,
				"socket_output": inputAddress,
				"api_method":    tool,
			},
		}, argv...)
	logdebug("Setting custom envvars to emulate Termux am wrapper script")
	// This block emulates the wrapper around `am` that's normally
	// used by termux.
	abiList, _ := am.Getprop("ro.product.cpu.abilist64")
	ldLibPath := "/system/lib" + abiList
	pathlessEnv := pathlessEnv()
	pathPlus := os.Getenv("PATH") + ":/system/bin"
	bc.Env = append(pathlessEnv, "LD_LIBRARY_PATH="+ldLibPath, "PATH="+pathPlus)
	logdebug("In execAMBroadcast, executing and awaiting finish", ctx{"inputAddress": inputAddress, "outputAddress": outputAddress, "argv": argv})
	op, err := bc.Output()
	logdebug("In execAMBroadcast, command output received", ctx{"output": string(op), "error": err})
	return err
}

// goroutine that transmits input buffer to a unix socket.
func transmitStdinToSocket(wg *sync.WaitGroup, ecbL func(error) bool, stdin io.Reader, unixSocket *net.UnixAddr) {
	defer wg.Done()
	logdebug("transmitStdinToSocket: start", ctx{"socketName": unixSocket.Name})
	outputClientSocket, err := net.Listen("unix", unixSocket.Name)
	if ecbL(err) {
		return
	}
	logdebug("transmitStdinToSocket: getting conn")
	conn, err := outputClientSocket.Accept()
	if ecbL(err) {
		return
	}
	logdebug("transmitStdinToSocket: copying stdin to conn")
	copiedbytes, err := io.Copy(conn, stdin)
	if ecbL(err) {
		return
	}
	logdebug("transmitStdinToSocket: copied successfully.", ctx{"copiedbytes": copiedbytes})
	logdebug("transmitSocketToStdout: closing conn")
	err = conn.Close()
	if ecbL(err) {
		return
	}
	logdebug("transmitStdinToSocket: finished")
}

// goroutine that transmits received data to an output buffer.
func transmitSocketToStdout(wg *sync.WaitGroup, ecbL func(error) bool, stdout io.Writer, unixSocket *net.UnixAddr) {
	defer wg.Done()
	logdebug("transmitSocketToStdout: start", ctx{"socketName": unixSocket.Name})
	inputClientSocket, err := net.Listen("unix", unixSocket.Name)
	if ecbL(err) {
		return
	}
	logdebug("transmitSocketToStdout: getting conn")
	conn, err := inputClientSocket.Accept()
	if ecbL(err) {
		return
	}
	logdebug("transmitSocketToStdout: copying conn output to output")
	_, err = io.Copy(stdout, conn)
	if ecbL(err) {
		return
	}
	logdebug("transmitSocketToStdout: Closing conn")
	err = conn.Close()
	if ecbL(err) {
		return
	}
	logdebug("transmitSocketToStdout: finished")
}

// Represents main() from the Termux-api C code.
// TODO: Needs serious clean-up.
func toolExec(stdin io.Reader, tool string, toolargs ...string) ([]byte, error) {
	// Prep bits
	wg := new(sync.WaitGroup)
	inputSocket := getUnixSocket()  // net.UnixAddr{Name: "@" + inputAddress, Net: "unix"}
	outputSocket := getUnixSocket() //net.UnixAddr{Name: "@" + outputAddress, Net: "unix"}
	var errs []error
	// Pitch off goroutines and wait.
	// Can this be a blocking call? Await proper test suite.
	// err := execAMBroadcast(inputSocket.Name, outputSocket.Name, tool, toolargs...)
	wg.Add(1)
	go func(wg *sync.WaitGroup, errcb func(error) bool) {
		defer wg.Done()
		err := execAMBroadcast(inputSocket.Name, outputSocket.Name, tool, toolargs...)
		if errcb(err) {
			return
		}
		loginfo("Finished broadcast goroutine without error")
	}(wg, makeErrCallback("broadCastGoroutine", errs))
	if stdin != nil {
		wg.Add(1)
		go transmitStdinToSocket(wg, makeErrCallback("transmitStdinToSocket", errs), stdin, &outputSocket) // Post stdin from this process to socket
	}
	stdoutBuf := new(bytes.Buffer)
	wg.Add(1)
	go transmitSocketToStdout(wg, makeErrCallback("transmitSocketToStdout", errs), stdoutBuf, &inputSocket) // Read from socket to this process' stdout
	logdebug("Awaiting goroutines to finish")
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return stdoutBuf.Bytes(), nil
}
