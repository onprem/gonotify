package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MakeRequest fires a new request to the Twilio API
func (t *Twilio) MakeRequest(path string, data url.Values) (map[string]interface{}, error) {
	endpoint := strings.Join([]string{t.BaseURL, t.APIVersion, "Accounts", t.AccountSID, path, ".json"}, "/")

	dataReader := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", endpoint, dataReader)
	req.SetBasicAuth(t.AccountSID, t.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twilio api, status: %s", resp.Status)
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response")
	}

	return result, nil
}
