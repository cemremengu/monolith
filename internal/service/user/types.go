package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user for profile and administration purposes.
type User struct {
	ID         uuid.UUID  `json:"id"               db:"id"`
	Username   string     `json:"username"         db:"username"`
	Email      string     `json:"email"            db:"email"`
	Name       *string    `json:"name"             db:"name"`
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
