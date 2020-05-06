package twilio

import "net/http"

// Twilio represent a twilio client
type Twilio struct {
	AccountSID string
	AuthToken  string
	BaseURL    string
	APIVersion string
	Client     *http.Client
}

// NewClient creates a new Twilio client
func NewClient(sid, token string) *Twilio {
	client := &http.Client{}
	return &Twilio{AccountSID: sid, AuthToken: token, BaseURL: "https://api.twilio.com", APIVersion: "2010-04-01", Client: client}
}

func (t *Twilio) bootstrapRequest(r *http.Request) {
	r.SetBasicAuth(t.AccountSID, t.AuthToken)
}
