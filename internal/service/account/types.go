package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID         uuid.UUID  `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Name       *string    `json:"name"`
	Password   string     `json:"-"`
	Avatar     *string    `json:"avatar,omitempty"`
	IsAdmin    bool       `json:"isAdmin"`
	Language   *string    `json:"language"`
	Theme      *string    `json:"theme"`
	Timezone   *string    `json:"timezone"`
	LastSeenAt *time.Time `json:"lastSeenAt"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"`
}

type UpdatePreferencesRequest struct {
	Language *string `json:"language"`
	Theme    *string `json:"theme"`
	Timezone *string `json:"timezone"`
}

type CreateAccountRequest struct {
	Username string  `json:"username" validate:"required"`
	Name     string  `json:"name"     validate:"required"`
	Email    string  `json:"email"    validate:"required,email"`
	Password *string `json:"password"`
	IsAdmin  *bool   `json:"isAdmin"`
}

type UpdateAccountRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

type InviteUsersRequest struct {
	Emails  []string `json:"emails" validate:"required,dive,email"`
	IsAdmin bool     `json:"isAdmin"`
}

type InviteUsersResponse struct {
	Success []Account           `json:"success"`
	Failed  []InviteUserFailure `json:"failed"`
}

type InviteUserFailure struct {
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword"     validate:"required,min=8"`
}
