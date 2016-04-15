package termux

import "bytes"

// Toast displays toast.
func Toast(content string, short bool) error {
	var args []string
	r := bytes.NewBuffer([]byte(content))
	if short {
		args = append(args, []string{"--ez", "short", "true"}...)
	}
	_, err := toolExec(r, "Toast", args...)
	return err
}
