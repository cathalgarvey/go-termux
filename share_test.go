package termux

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateShareArgs(t *testing.T) {
	for _, valid := range []string{"edit", "send", "view"} {
		err := validateShareArguments(valid)
		assert.Nil(t, err)
	}
	for _, invalid := range []string{"make", "del", "Edit", "send ", "VIEW"} {
		err := validateShareArguments(invalid)
		assert.EqualError(t, err, "Bad share action; must be (edit|send|view)")
	}
}

func Test_ShareConstruction(t *testing.T) {
	// share(execF toolExecFunc, title, action, contentType string, defaultAction bool, content io.Reader) error {
	ef := prepDummyExecFunc(t, []byte("foo"), "Share", []string{"--es", "title", "title.bar", "--es", "action", "view", "--es", "content-type", "text/idiocy"}, []byte{})
	err := share(ef, "title.bar", "view", "text/idiocy", false, bytes.NewReader([]byte("foo")))
	assert.Nil(t, err)
}
