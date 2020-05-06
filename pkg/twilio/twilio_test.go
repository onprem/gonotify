package twilio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// MockTwilio return a mocked version of Twilio client
func MockTwilio(code int, response string) (*Twilio, *httptest.Server) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authH := r.Header.Get("Authorization")
		if authH != "Basic YWJjOjEyMw==" { // base64 for abc:123
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(code)
		fmt.Fprintln(w, response)
	}))

	c := ts.Client()

	return &Twilio{
		AccountSID: "abc",
		AuthToken:  "123",
		BaseURL:    ts.URL,
		APIVersion: "2010-04-01",
		Client:     c,
	}, ts
}
