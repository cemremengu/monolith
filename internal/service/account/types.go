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
	IsDisabled bool       `json:"isDisabled"`
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
