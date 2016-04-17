package termux

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

var (
	// ErrNoRecipientNumber is returned if the number string is empty
	ErrNoRecipientNumber = errors.New("No recipient number specified for SMS")
)

// SMS is an SMS message from SmsInbox
// TODO need struct format, this API doesn't work when one uses SMSSecure :)
type SMS struct {
	Read     bool
	Sender   string // Presumably only if contact is present
	Number   string
	Received string // YYYY-MM-DD HH:MM
	Body     string
}

// SMSSend sends an SMS
func SMSSend(content, number string) error {
	if number == "" {
		return ErrNoRecipientNumber
	}
	_, err := toolExec(bytes.NewBuffer([]byte(content)), "SmsSend", "--es", "recipient", number)
	return err
}

// SMSInbox returns the main messaging app's Inbox. It returns nothing if the
// phone is set to use SMSSecure or other private SMS apps.
// Limit is how many to return, most recent first.
// Offset is how many to skip back from most recent before counting out <Limit> smses.
// ShowDates and ShowPhoneNumbers appear to be the default anyway, so have been elided.
func SMSInbox(limit, offset int) ([]SMS, error) {
	var (
		resp []SMS
		args []string
	)
	args = append(args, []string{"--ez", "show-dates", "true"}...)
	args = append(args, []string{"--ez", "show-phone-numbers", "true"}...)
	if offset != 0 {
		args = append(args, []string{"--ei", "offset", strconv.Itoa(offset)}...)
	}
	if limit != 0 {
		args = append(args, []string{"--ei", "limit", strconv.Itoa(limit)}...)
	}
	bytesResp, err := toolExec(nil, "SmsInbox", args...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytesResp, &resp)
	return resp, err
}

// Time returns a structured time representing the received message timestamp.
func (sms SMS) Time() (time.Time, error) {
	// Reftime: Mon Jan 2 15:04:05 -0700 MST 2006
	return time.ParseInLocation("2006-01-02 15:04", sms.Received, time.Local)
}
