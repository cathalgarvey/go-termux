package termux

// ClipboardGet returns the contents of the Clipboard
func ClipboardGet() (string, error) {
	return clipboardGet(toolExec)
}

func clipboardGet(execF toolExecFunc) (string, error) {
	strv, err := execF(nil, "Clipboard")
	return string(strv), err
}

// ClipboardSet sets the contents of the Clipboard
func ClipboardSet(val string) error {
	return clipboardSet(toolExec, val)
}

func clipboardSet(execF toolExecFunc, val string) error {
	_, err := execF(nil, "Clipboard", "--es", "text", val)
	return err
}
