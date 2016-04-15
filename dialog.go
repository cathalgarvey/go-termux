package termux

import "bytes"

// Dialog displays an input dialog and returns the input data.
func Dialog(title, hint string, multiline, passwordInput bool) (string, error) {
	var args []string
	if title != "" {
		args = append(args, []string{"--es", "input_title", title}...)
	}
	if hint != "" {
		args = append(args, []string{"--es", "input_hint", hint}...)
	}
	if multiline {
		args = append(args, []string{"--ez", "multiple_lines", "true"}...)
	}
	if passwordInput {
		args = append(args, []string{"--es", "input_type", "password"}...)
	}
	conbytes, err := toolExec(bytes.NewBuffer(nil), "Dialog", args...)
	return string(conbytes), err
}
