package termux

import "encoding/json"

// Contact is the stripped down data returned by termux-contact-list
type Contact struct {
	Name   string
	Number string
}

// ContactList returns a list of Contacts. These appear to be stripped down to
// only the name and a phone number.
func ContactList() ([]Contact, error) {
	return contactList(toolExec)
}

func contactList(execF toolExecFunc) ([]Contact, error) {
	var cons []Contact
	conbytes, err := execF(nil, "ContactList")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(conbytes, &cons)
	return cons, err
}
