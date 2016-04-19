package termux

import "errors"

var (
	// ErrNoOutputFile is returned from the CameraPhoto API if an output filename
	// is expected and not given.
	ErrNoOutputFile = errors.New("Must provide an output filename")

	// ErrNoCameraSpecified is returned from the CameraPhoto API if cameraID is blank.
	ErrNoCameraSpecified = errors.New("Must specify a camera ID")

	// ErrNoURLSpecified is returned from the Download API if URL==""
	ErrNoURLSpecified = errors.New("No URL specified for download")

	// ErrBadLocationRequest is returned from the Location API if request param is not (once|last|updates)
	ErrBadLocationRequest = errors.New("Bad Location `request` parameter")

	// ErrBadLocationProvider is returned from the Location API if provider param is not (gps|network|passive)
	ErrBadLocationProvider = errors.New("Bad Location `provider` parameter")

	// ErrUpdatesNotSupportedYet is returned from the Location API for "updates" because it's not supported yet.
	ErrUpdatesNotSupportedYet = errors.New("'updates' mode not yet supported")

	// ErrNoNotificationArgsProvided is returned from the Notification API if neither title or content is given
	ErrNoNotificationArgsProvided = errors.New("Either title or content must be provided")

	// ErrBadShareAction is returned from the Share API if 'action' is a bad value
	ErrBadShareAction = errors.New("Bad share action; must be (edit|send|view)")

	// ErrNoRecipientNumber is returned from the SMS API if the number string is empty
	ErrNoRecipientNumber = errors.New("No recipient number specified for SMS")

	// ErrNoTTSEngineProvided is returned from the TTS API if there's no engine specified
	ErrNoTTSEngineProvided = errors.New("No engine specified for TTS")

	// ErrNoTTSContentProvided is returned from the TTS API if no content is given to speak
	ErrNoTTSContentProvided = errors.New("No content given to TTS engine")
)
