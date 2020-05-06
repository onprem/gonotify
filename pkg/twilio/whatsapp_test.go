package twilio

import (
	"net/http"
	"testing"
)

func TestSendWhatsAppSuccess(t *testing.T) {
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "delivered" }`)
	defer ts.Close()

	err := tw.SendWhatsApp("whatsapp:+1789789789", "whatsapp:+9112312312", "test")
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}
}

func TestSendWhatsAppFailure(t *testing.T) {
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "failed", "error_message": "unknown error", "error_code": 30008 }`)
	defer ts.Close()

	err := tw.SendWhatsApp("whatsapp:+1789789789", "whatsapp:+9112312312", "test")
	if err == nil {
		t.Error("Expected error but got success")
	}
}

func TestSendWhatsAppInvalidNumber(t *testing.T) {
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "failed", "error_message": "unknown error", "error_code": 30008 }`)
	defer ts.Close()

	err := tw.SendWhatsApp("whats:+1789789789", "whatsapp:+9112312312", "test")
	if err != ErrInvalidWhatsAppNumber {
		t.Errorf("Expected error '%s' but got success", ErrInvalidWhatsAppNumber)
	}

	err = tw.SendWhatsApp("whatsapp:+1789789789", "+9112312312", "test")
	if err != ErrInvalidWhatsAppNumber {
		t.Errorf("Expected error '%s' but got success", ErrInvalidWhatsAppNumber)
	}
}

func TestSendWhatsAppEmptyBody(t *testing.T) {
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "failed", "error_message": "unknown error", "error_code": 30008 }`)
	defer ts.Close()

	err := tw.SendWhatsApp("whatsapp:+1789789789", "whatsapp:+9112312312", "")
	if err != ErrEmptyWhatsAppBody {
		t.Errorf("Expected error '%s' but got success", ErrEmptyWhatsAppBody)
	}
}
