package login

import (
	"liquide_assignment/internal/config"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Pan      string `json:"pan"`
	Mobile   string `json:"mobile"`
	Name     string `json:"name"`
	UserId   string `json:"-"`
}

type Response struct {
	Message   string           `json:"message"`
	Status    config.Status    `json:"status"`
	ErrorCode config.ErrorCode `json:"error_code"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"device_id,omitempty"`
	Platform string `json:"platform,omitempty"`
}

type UserSessionDetails struct {
	ExpiryTime time.Time `json:"expiry_time"`
	IsBlocked  bool      `json:"is_blocked"`
}

func (r *LoginRequest) SetDefaults() {
	if r.Platform == "" {
		r.Platform = "WEB"
	}
}

type tokenWithExp struct {
	Token string
	Exp   time.Time
}

var panRegex = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)

func (s *SignupRequest) Validate() bool {
	if s.Email == "" {
		return false
	}

	_, err := mail.ParseAddress(s.Email)
	if err != nil {
		return false
	}

	if s.Password == "" {
		return false
	}

	if len(s.Password) < 8 {
		return false
	}

	if s.Name == "" {
		return false
	}

	if s.Mobile == "" {
		return false
	}

	if len(s.Mobile) != 10 {
		return false
	}

	if s.Pan == "" {
		return false
	}

	if len(s.Pan) != 10 {
		return false
	}

	if !panRegex.MatchString(strings.ToUpper(strings.TrimSpace(s.Pan))) {
		return false
	}

	return true
}
