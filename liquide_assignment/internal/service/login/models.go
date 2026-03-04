package login

import (
	"liquide_assignment/internal/config"
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
