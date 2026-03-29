package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"monolith/internal/database"
	"monolith/internal/database/dbsqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db      *database.DB
	queries *dbsqlc.Queries
}

func NewService(db *database.DB) *Service {
	return &Service{
		db:      db,
		queries: db.Queries(),
	}
}

func (s *Service) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) UserExists(ctx context.Context, email, username string) (bool, error) {
	_, err := s.queries.UserExists(ctx, dbsqlc.UserExistsParams{
		Email:    email,
		Username: username,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*Account, error) {
	if req.Password == "" || len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	exists, err := s.UserExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	row, err := s.queries.WithTx(tx).RegisterAccount(ctx, dbsqlc.RegisterAccountParams{
		Username: req.Username,
		Email:    req.Email,
		Name:     textFromString(req.Name),
		Password: textFromStringPtr(&hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	account := accountFromRegisterRow(row)
	return &account, nil
}

func (s *Service) GetAccountByLogin(ctx context.Context, login string) (*Account, error) {
	row, err := s.queries.GetAccountByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		Password:   stringFromText(row.Password),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
	return &account, nil
}

func (s *Service) GetAccountByID(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	row, err := s.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
	return &account, nil
}

func (s *Service) UpdateLastSeen(ctx context.Context, accountID uuid.UUID) error {
	return s.queries.UpdateLastSeen(ctx, accountID)
}

func (s *Service) UpdatePreferences(
	ctx context.Context,
	accountID uuid.UUID,
	req UpdatePreferencesRequest,
) (*Account, error) {
	row, err := s.queries.UpdatePreferences(ctx, dbsqlc.UpdatePreferencesParams{
		Language: textFromStringPtr(req.Language),
		Theme:    textFromStringPtr(req.Theme),
		Timezone: textFromStringPtr(req.Timezone),
		ID:       accountID,
	})
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
	return &account, nil
}

func (s *Service) GetAccounts(ctx context.Context) ([]Account, error) {
	rows, err := s.queries.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make([]Account, len(rows))
	for i, row := range rows {
		accounts[i] = Account{
			ID:         row.ID,
			Username:   row.Username,
			Email:      row.Email,
			Name:       stringPtrFromText(row.Name),
			Avatar:     stringPtrFromText(row.Avatar),
			IsAdmin:    boolFromPgBool(row.IsAdmin),
			Language:   stringPtrFromText(row.Language),
			Theme:      stringPtrFromText(row.Theme),
			Timezone:   stringPtrFromText(row.Timezone),
			LastSeenAt: row.LastSeenAt,
			Status:     row.Status,
			CreatedAt:  timeFromPtr(row.CreatedAt),
			UpdatedAt:  timeFromPtr(row.UpdatedAt),
		}
	}
	return accounts, nil
}

func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (*Account, error) {
	row, err := s.queries.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		Avatar:     stringPtrFromText(row.Avatar),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
	return &account, nil
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	var hashedPassword *string
	status := "pending"
	isAdmin := false

	if req.Password != nil && *req.Password != "" {
		hashed, err := hashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = &hashed
		status = "active"
	}

	if req.IsAdmin != nil {
		isAdmin = *req.IsAdmin
	}

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	row, err := s.queries.WithTx(tx).CreateAccount(ctx, dbsqlc.CreateAccountParams{
		Username: req.Username,
		Name:     textFromString(req.Name),
		Email:    req.Email,
		Password: textFromStringPtr(hashedPassword),
		IsAdmin:  pgtype.Bool{Bool: isAdmin, Valid: true},
		Status:   status,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	account := accountFromCreateRow(row)
	return &account, nil
}

func (s *Service) InviteUsers(ctx context.Context, req InviteUsersRequest) (*InviteUsersResponse, error) {
	response := &InviteUsersResponse{
		Success: []Account{},
		Failed:  []InviteUserFailure{},
	}

	for _, email := range req.Emails {
		username := deriveUsernameFromEmail(email)

		_, err := s.queries.UserExists(ctx, dbsqlc.UserExistsParams{
			Email:    email,
			Username: username,
		})

		if err == nil {
			response.Failed = append(response.Failed, InviteUserFailure{
				Email:  email,
				Reason: "User already exists",
			})
			continue
		}

		row, err := s.queries.CreateInvitedAccount(ctx, dbsqlc.CreateInvitedAccountParams{
			Username: username,
			Email:    email,
			Name:     textFromString(username),
			IsAdmin:  pgtype.Bool{Bool: req.IsAdmin, Valid: true},
		})
		if err != nil {
			response.Failed = append(response.Failed, InviteUserFailure{
				Email:  email,
				Reason: fmt.Sprintf("Failed to create user: %v", err),
			})
			continue
		}

		response.Success = append(response.Success, accountFromInvitedRow(row))
	}

	return response, nil
}

func (s *Service) UpdateAccount(ctx context.Context, id uuid.UUID, req UpdateAccountRequest) (*Account, error) {
	row, err := s.queries.UpdateAccount(ctx, dbsqlc.UpdateAccountParams{
		Username: req.Username,
		Name:     textFromString(req.Name),
		Email:    req.Email,
		ID:       id,
	})
	if err != nil {
		return nil, err
	}

	account := Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		Avatar:     stringPtrFromText(row.Avatar),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
	return &account, nil
}

func (s *Service) DisableAccount(ctx context.Context, id uuid.UUID) error {
	return s.queries.DisableAccount(ctx, id)
}

func (s *Service) EnableAccount(ctx context.Context, id uuid.UUID) error {
	return s.queries.EnableAccount(ctx, id)
}

// TODO: this should trigger deleting other stuff
func (s *Service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteAccount(ctx, id)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func deriveUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		username := strings.ToLower(parts[0])
		username = strings.ReplaceAll(username, ".", "")
		username = strings.ReplaceAll(username, "+", "")
		return username
	}
	return email
}

func (s *Service) ChangePassword(ctx context.Context, accountID uuid.UUID, req ChangePasswordRequest) error {
	if req.NewPassword == "" || len(req.NewPassword) < 8 {
		return ErrPasswordTooShort
	}

	row, err := s.queries.GetAccountByIDWithPassword(ctx, accountID)
	if err != nil {
		return err
	}

	if err := s.ValidatePassword(stringFromText(row.Password), req.CurrentPassword); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.queries.UpdatePassword(ctx, dbsqlc.UpdatePasswordParams{
		Password: textFromStringPtr(&hashedPassword),
		ID:       accountID,
	})
}

func textFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func textFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func stringPtrFromText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func stringFromText(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func boolFromPgBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func timeFromPtr(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func accountFromRegisterRow(row dbsqlc.RegisterAccountRow) Account {
	return Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
}

func accountFromCreateRow(row dbsqlc.CreateAccountRow) Account {
	return Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		Avatar:     stringPtrFromText(row.Avatar),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
}

func accountFromInvitedRow(row dbsqlc.CreateInvitedAccountRow) Account {
	return Account{
		ID:         row.ID,
		Username:   row.Username,
		Email:      row.Email,
		Name:       stringPtrFromText(row.Name),
		Avatar:     stringPtrFromText(row.Avatar),
		IsAdmin:    boolFromPgBool(row.IsAdmin),
		Language:   stringPtrFromText(row.Language),
		Theme:      stringPtrFromText(row.Theme),
		Timezone:   stringPtrFromText(row.Timezone),
		LastSeenAt: row.LastSeenAt,
		Status:     row.Status,
		CreatedAt:  timeFromPtr(row.CreatedAt),
		UpdatedAt:  timeFromPtr(row.UpdatedAt),
	}
}
