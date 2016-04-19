package termux

// Download issues a download using the system download manager; i.e. with a
// title and description in the status bar.
func Download(URL, title, description string) error {
	return download(toolExec, URL, title, description)
}

func download(execF toolExecFunc, URL, title, description string) error {
	var args []string
	if URL == "" {
		return ErrNoURLSpecified
	}
	if title != "" {
		args = append(args, []string{"--es", "title", title}...)
	}
	if description != "" {
		args = append(args, []string{"--es", "description", description}...)
	}
	args = append(args, URL)
	_, err := execF(nil, "Download", args...)
	return err
}
