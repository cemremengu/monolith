package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user for profile and administration purposes.
type User struct {
	ID         uuid.UUID  `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Name       *string    `json:"name"`
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

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}
