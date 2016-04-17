package termux

import "errors"

var (
	// ErrNoURLSpecified is returned if URL=="" in Download
	ErrNoURLSpecified = errors.New("No URL specified for download")
)

// Download issues a download using the system download manager; i.e. with a
// title and description in the status bar.
func Download(URL, title, description string) error {
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
	_, err := toolExec(nil, "Download", args...)
	return err
}
