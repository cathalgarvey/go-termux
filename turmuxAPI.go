/*Package termux provides a clean Go API for making calls to Termux's Android API tools.

It attempts to reimplement and replace the C code from termux-api, and to provide
a clean functional interface to all the api script wrappers, so that development
of functional termux "apps" can be as easy as:

1. Install Termux
2. Install Termux:API (just the app: no need to install termux-api!)
3. Write and test your go app from Termux shell (possibly even compiling from Termux!)
4. Optional: Use Termux:Widget to create Android desktop shortcuts for your new
   termux/Go "app"
*/
package termux

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"

	am "github.com/cathalgarvey/androidam"
	"github.com/twinj/uuid"
)

var (
	// ErrZeroLengthResponse is returned if the response is (unexpectedly) zero in length
	ErrZeroLengthResponse = errors.New("Zero length response received")

	// ErrNoOutputFile is returned if an output filename is expected
	ErrNoOutputFile = errors.New("Must provide an output filename")
)

// execAMBroadcast is a rewrite of exec_am_broadcast from:
// https://github.com/termux/termux-packages/blob/master/packages/termux-api/termux-api.c#L27
func execAMBroadcast(inputAddress, outputAddress string, tool string, argv ...string) error {
	logdebug("In execAMBroadcast, prepping call", ctx{"inputAddress": inputAddress, "outputAddress": outputAddress, "argv": argv})
	bc := am.Broadcast(
		&am.Opts{
			AMPath: "/data/data/com.termux/files/usr/bin/am",
		},
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
	logdebug("In execAMBroadcast, executing and awaiting finish", ctx{"inputAddress": inputAddress, "outputAddress": outputAddress, "argv": argv})
	op, err := bc.Output()
	logdebug("In execAMBroadcast, command output received", ctx{"output": string(op), "error": err})
	return err
}

func generateUUID() string {
	return uuid.NewV4().String()
}

func transmitStdinToSocket(wg *sync.WaitGroup, echan chan error, stdin io.Reader, unixSocket *net.UnixAddr) {
	defer wg.Done()
	logdebug("transmitStdinToSocket: start", ctx{"socketName": unixSocket.Name})
	outputClientSocket, err := net.Listen("unix", unixSocket.Name)
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitStdinToSocket", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitStdinToSocket: getting conn")
	conn, err := outputClientSocket.Accept()
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitStdinToSocket", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitStdinToSocket: copying stdin to conn")
	copiedbytes, err := io.Copy(conn, stdin)
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitStdinToSocket", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitStdinToSocket: copied successfully.", ctx{"copiedbytes": copiedbytes})
	logdebug("transmitSocketToStdout: closing conn")
	err = conn.Close()
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitStdinToSocket", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitStdinToSocket: finished")
}

func transmitSocketToStdout(wg *sync.WaitGroup, echan chan error, stdout io.Writer, unixSocket *net.UnixAddr) {
	defer wg.Done()
	logdebug("transmitSocketToStdout: start", ctx{"socketName": unixSocket.Name})
	inputClientSocket, err := net.Listen("unix", unixSocket.Name)
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitSocketToStdout", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitSocketToStdout: getting conn")
	conn, err := inputClientSocket.Accept()
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitSocketToStdout", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitSocketToStdout: copying conn output to output")
	_, err = io.Copy(stdout, conn)
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitSocketToStdout", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitSocketToStdout: Closing conn")
	err = conn.Close()
	if err != nil {
		logerror("Error in socket thread", ctx{"function": "transmitSocketToStdout", "error": err.Error()})
		echan <- err
		return
	}
	logdebug("transmitSocketToStdout: finished")
}

// toolExecFunc represents the interface all of the other tools construct args
// for and pass-through. It's interfaced here to allow for testing.
type toolExecFunc func(stdin io.Reader, tool string, toolargs ...string) ([]byte, error)

// Represents main() from the Termux-api C code.
func toolExec(stdin io.Reader, tool string, toolargs ...string) ([]byte, error) {
	// Prep bits
	stdoutBuf := new(bytes.Buffer)
	inputAddress := generateUUID()  // "This Program reads from it"
	outputAddress := generateUUID() // "This Program writes to it"
	inputSocket := net.UnixAddr{Name: "@" + inputAddress, Net: "unix"}
	outputSocket := net.UnixAddr{Name: "@" + outputAddress, Net: "unix"}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	var errs []error
	ech := make(chan error, 3) // Capacity for up to three errs, should be all it needs.
	loginfo("Starting broadcast goroutine", ctx{
		"inputAddress":  inputAddress,
		"outputAddress": outputAddress,
	})
	// Pitch off three goroutines and wait.
	go func(wg *sync.WaitGroup, errchan chan error) {
		defer wg.Done()
		err := execAMBroadcast(inputAddress, outputAddress, tool, toolargs...)
		if err != nil {
			logerror("Error in broadcast goroutine", ctx{"error": err})
			errchan <- err
			return
		}
		loginfo("Finished broadcast goroutine without error")
	}(wg, ech)
	if stdin != nil {
		wg.Add(1)
		go transmitStdinToSocket(wg, ech, stdin, &outputSocket) // Post stdin from this process to socket
	}
	// TODO: If either of the above fail, will this one wait forever?
	// Should it use a select statement and a timeout?
	go transmitSocketToStdout(wg, ech, stdoutBuf, &inputSocket) // Read from socket to this process' stdout
	logdebug("Awaiting goroutines to finish")
	wg.Wait()
	close(ech)
	logdebug("Goroutine waitgroup finished, cleaning up")
	// Clean up, get data, get out
	for err := range ech {
		logerror("Error in socket threads", ctx{"err": err})
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		logerror("Errors occurred, aborting", ctx{"errors": errs})
		return nil, errs[0]
	}
	logdebug("Success, returning buffer")
	return stdoutBuf.Bytes(), nil
}
