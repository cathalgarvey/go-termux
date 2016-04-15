package termux

import (
	"bytes"
	"errors"
)

var (
	// ErrNoNotificationArgsProvided is returned if neither title or content is given
	ErrNoNotificationArgsProvided = errors.New("Either title or content must be provided")
)

// Notification posts a notification to the system tray
func Notification(title, content, id, uri string) error {
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
	_, err := toolExec(bytes.NewBuffer(nil), "Notification", args...)
	return err
}
