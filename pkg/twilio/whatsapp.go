package twilio

import (
	"fmt"
	"net/url"
	"strings"
)

var messagePath = "Messages"

// SendWhatsApp sends an message over WhatsApp
func (t *Twilio) SendWhatsApp(from, to, body string) error {
	wPrefix := "whatsapp:"
	if !strings.HasPrefix(from, wPrefix) || !strings.HasPrefix(to, wPrefix) {
		return ErrInvalidWhatsAppNumber
	}
	if body == "" {
		return ErrEmptyWhatsAppBody
	}

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", from)
	data.Set("Body", body)

	resp, err := t.MakeRequest(messagePath, data)
	if err != nil {
		return err
	}

	if resp["status"] == "failed" || resp["status"] == "undelivered" {
		return fmt.Errorf("msg: %s code: %s", resp["error_message"], resp["error_code"])
	}

	return nil
}

// ErrInvalidWhatsAppNumber is thrown when number is not suffixed with whatsapp:
var ErrInvalidWhatsAppNumber = fmt.Errorf("Invalid WhatsApp number")

// ErrEmptyWhatsAppBody is thrown when body parameter is empty
var ErrEmptyWhatsAppBody = fmt.Errorf("Message Body is empty")
