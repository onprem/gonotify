package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

	return strings.ReplaceAll(phonenumbers.Format(ph, phonenumbers.INTERNATIONAL), " ", ""), nil
}
