package termux

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// toolExecFunc func(stdin io.Reader, tool string, toolargs ...string) ([]byte, error)
// prepDummyExecFunc creates a callback function that validates
//  the expected inputs and fails tests if any don't match.
// It returns the expected output if no error occurs.
func prepDummyExecFunc(t *testing.T, expectedStdin []byte, expectedTool string, expectedArgs []string, expectedOutput []byte) toolExecFunc {
	return func(stdin io.Reader, tool string, toolargs ...string) ([]byte, error) {
		assert.EqualValues(t, expectedTool, tool)
		assert.EqualValues(t, expectedArgs, toolargs)
		stdinput, err := ioutil.ReadAll(stdin)
		assert.Nil(t, err)
		assert.EqualValues(t, expectedStdin, stdinput)
		// To assist in testing..
		return expectedOutput, nil
	}
}
