package termux

import "bytes"

// ClipboardGet returns the contents of the Clipboard
func ClipboardGet() (string, error) {
	strv, err := toolExec(bytes.NewBuffer(nil), "Clipboard")
	return string(strv), err
}

// ClipboardSet sets the contents of the Clipboard
func ClipboardSet(val string) error {
	_, err := toolExec(bytes.NewBuffer(nil), "Clipboard", "--es", "text", val)
	return err
}
