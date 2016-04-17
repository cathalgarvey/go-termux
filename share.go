package termux

import (
	"errors"
	"io"
)

var (
	// ErrBadShareAction is returned if 'action' is a bad value
	ErrBadShareAction = errors.New("Bad share action; must be (edit|send|view)")
)

// Share sends a file using a share dialog or default share action
func Share(title, action, contentType string, defaultAction bool, content io.Reader) error {
	return share(toolExec, title, action, contentType, defaultAction, content)
}

func share(execF toolExecFunc, title, action, contentType string, defaultAction bool, content io.Reader) error {
	if err := validateShareArguments(action); err != nil {
		return err
	}
	var args []string
	if title != "" {
		args = append(args, []string{"--es", "title", title}...)
	}
	args = append(args, []string{"--es", "action", action}...)
	if contentType != "" {
		args = append(args, []string{"--es", "content-type", contentType}...)
	} // else { guessContentType }
	if defaultAction {
		args = append(args, []string{"--ez", "default-receiver", "true"}...)
	}
	_, err := execF(content, "Share", args...)
	return err
}

func validateShareArguments(action string) error {
	switch action {
	case "edit", "send", "view":
		return nil
	default:
		return ErrBadShareAction
	}
}
