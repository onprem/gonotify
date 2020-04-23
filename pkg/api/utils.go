package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/nyaruka/phonenumbers"
)

func getVerificationCode() string {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(899999) + 100000
	return strconv.Itoa(n)
}

func parsePhone(phone string) (string, error) {
	ph, err := phonenumbers.Parse(phone, "IN")
	if err != nil || !phonenumbers.IsValidNumber(ph) {
		return "", fmt.Errorf("invalid phone number")
	}

	return phonenumbers.Format(ph, phonenumbers.E164), nil
}
