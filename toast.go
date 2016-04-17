package termux

import "bytes"

// Toast displays toast.
func Toast(content string, short bool) error {
	var args []string
	if short {
		args = append(args, []string{"--ez", "short", "true"}...)
	}
	_, err := toolExec(bytes.NewBuffer([]byte(content)), "Toast", args...)
	return err
}
