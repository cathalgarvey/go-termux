package termux

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

var (
	// ErrNoRecipientNumber is returned if the number string is empty
	ErrNoRecipientNumber = errors.New("No recipient number specified for SMS")
)

// SMS is an SMS message from SmsInbox
// TODO need struct format, this API doesn't work when one uses SMSSecure :)
type SMS struct {
	Read     bool
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

// SMSInbox returns the SMSInbox
func SMSInbox(limit, offset int, showDates, showPhoneNumbers bool) ([]SMS, error) {
	var (
		resp []SMS
		args []string
	)
	if showDates {
		args = append(args, []string{"--ez", "show-dates", "true"}...)
	}
	if showPhoneNumbers {
		args = append(args, []string{"--ez", "show-phone-numbers", "true"}...)
	}
	if offset != 0 {
		args = append(args, []string{"--ei", "offset", strconv.Itoa(offset)}...)
	}
	if limit != 0 {
		args = append(args, []string{"--ei", "limit", strconv.Itoa(limit)}...)
	}
	bytesResp, err := toolExec(bytes.NewBuffer(nil), "SmsInbox", args...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytesResp, &resp)
	return resp, err
}
