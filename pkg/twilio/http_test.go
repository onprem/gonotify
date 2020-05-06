package twilio

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMakeRequestUnauthorized(t *testing.T) {
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "delivered" }`)
	defer ts.Close()
	tw.AccountSID = "def"

	data := url.Values{}
	data.Set("To", "+911212121221")
	data.Set("From", "+178978945")
	data.Set("Body", "test")

	res, err := tw.MakeRequest("Messages", data)
	if err == nil {
		t.Errorf("Expected Unauthorized, go success: %v", res)
	}
}
