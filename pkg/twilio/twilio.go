package twilio

import "net/http"

// Twilio represent a twilio client
type Twilio struct {
	AccountSID string
	AuthToken  string
	client     *http.Client
}

// NewClient creates a new Twilio client
func NewClient(sid, token string) *Twilio {
	client := &http.Client{}
	return &Twilio{AccountSID: sid, AuthToken: token, client: client}
}

func (t *Twilio) bootstrapRequest(r *http.Request) {
	r.SetBasicAuth(t.AccountSID, t.AuthToken)

}
