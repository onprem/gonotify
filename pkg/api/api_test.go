package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prmsrswt/gonotify/pkg/twilio"
)

var testConf = Config{
	TwilioSID:      "abc",
	TwilioToken:    "123",
	JWTSecret:      []byte("test"),
	WhatsAppFrom:   "whatsapp:+9111111111",
	WebHookAccount: gin.Accounts{"user": "pass"},
	VerifyTmpl:     template.Must(template.New("verify").Parse("Verification code is {{ .Code }}")),
	NotifTmpl:      template.Must(template.New("notif").Parse("You got {{ .Total }} new notifications.")),
}

// MockTwilio return a mocked version of Twilio client
func MockTwilio(code int, response string) (*twilio.Twilio, *httptest.Server) {
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

	return &twilio.Twilio{
		AccountSID: "abc",
		AuthToken:  "123",
		BaseURL:    ts.URL,
		APIVersion: "2010-04-01",
		Client:     c,
	}, ts
}

func testGin() (*API, *httptest.Server) {
	db, _ := sql.Open("sqlite3", ":memory:")
	tw, ts := MockTwilio(http.StatusOK, `{ "status": "delivered" }`)

	gin.SetMode(gin.ReleaseMode)
	g := &API{
		conf:         testConf,
		Gin:          gin.Default(),
		DB:           db,
		TwilioClient: tw,
		logger:       log.NewNopLogger(),
	}
	_ = bootstrapDB(db)
	g.Register()

	return g, ts
}

func testJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  1,
		"exp": time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenStr, _ := token.SignedString(testConf.JWTSecret)
	return tokenStr
}

func TestPing(t *testing.T) {
	router, ts := testGin()
	defer ts.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	router.Gin.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Unauthenticated request didn't got 401, expected 401, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+testJWT())
	router.Gin.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Request to /ping failed, expected 200, got %d", w.Code)
	}
}
