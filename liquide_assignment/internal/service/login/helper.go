package login

import (
	"crypto/rand"
	"encoding/hex"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var panRegex = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)

func checkUserValidity(request *SignupRequest) bool {
	if request.Email == "" {
		return false
	}

	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		return false
	}

	if request.Password == "" {
		return false
	}

	if len(request.Password) < 8 {
		return false
	}

	if request.Name == "" {
		return false
	}

	if request.Mobile == "" {
		return false
	}

	if len(request.Mobile) != 10 {
		return false
	}

	if request.Pan == "" {
		return false
	}

	if len(request.Pan) != 10 {
		return false
	}

	if !panRegex.MatchString(strings.ToUpper(strings.TrimSpace(request.Pan))) {
		return false
	}

	return true
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func generateReqToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
