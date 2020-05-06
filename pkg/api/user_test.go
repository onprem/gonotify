package api

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"encoding/json"
)

func TestUserRegister(t *testing.T) {
	router, ts := testGin()
	defer ts.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"/api/v1/register",
		bytes.NewBufferString(`{ "name": "Test", "phone": "+919955555555", "password": "123456"}`),
	)
	router.Gin.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("User registration failed; expected 200, got %d", w.Code)
	}

	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)

	errMsg, ok := res["error"]

	if ok {
		t.Errorf("User registration failed; expected success, got %s", errMsg)
	}
}
