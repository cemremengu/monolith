package auth

import (
	"time"
)

type UserAccount struct {
	ID         string     `json:"id"               db:"id"`
	Username   string     `json:"username"         db:"username"`
	Email      string     `json:"email"            db:"email"`
	Name       *string    `json:"name"             db:"name"`
	Password   string     `json:"-"                db:"password"`
	Avatar     *string    `json:"avatar,omitempty" db:"avatar"`
	IsAdmin    bool       `json:"isAdmin"          db:"is_admin"`
	Language   *string    `json:"language"         db:"language"`
	Theme      *string    `json:"theme"            db:"theme"`
	Timezone   *string    `json:"timezone"         db:"timezone"`
	LastSeenAt *time.Time `json:"lastSeenAt"       db:"last_seen_at"`
	IsDisabled bool       `json:"isDisabled"       db:"is_disabled"`
	CreatedAt  time.Time  `json:"createdAt"        db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt"        db:"updated_at"`
}

type LoginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserAccount `json:"user"`
}
