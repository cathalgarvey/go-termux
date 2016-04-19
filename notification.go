package termux

// Notification posts a notification to the system tray
func Notification(title, content, id, uri string) error {
	return notification(toolExec, title, content, id, uri)
}

func notification(execF toolExecFunc, title, content, id, uri string) error {
	var args []string
	if title == "" && content == "" {
		return ErrNoNotificationArgsProvided
	}
	if title != "" {
		args = append(args, []string{"--es", "title", title}...)
	}
	if content != "" {
		args = append(args, []string{"--es", "content", content}...)
	}
	if id != "" {
		args = append(args, []string{"--es", "id", id}...)
	}
	if uri != "" {
		args = append(args, []string{"--es", "url", uri}...)
	}
	_, err := execF(nil, "Notification", args...)
	return err
}
