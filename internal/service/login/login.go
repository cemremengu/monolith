package login

import (
	"context"
	"errors"
	"log/slog"

	"monolith/internal/database"
	"monolith/internal/service/account"
	"monolith/internal/service/ldap"
)

type Service struct {
	db             *database.DB
	accountService *account.Service
	ldapService    *ldap.Service
}

func NewService(db *database.DB, accountService *account.Service, ldapService *ldap.Service) *Service {
	return &Service{
		db:             db,
		accountService: accountService,
		ldapService:    ldapService,
	}
}

func (s *Service) Login(ctx context.Context, req UserLoginRequest) (*account.Account, error) {
	if s.ldapService != nil && s.ldapService.Enabled() {
		acc, err := s.loginLDAP(ctx, req)
		if err == nil {
			return acc, nil
		}

		if !errors.Is(err, ldap.ErrUserNotFound) && !errors.Is(err, ldap.ErrInvalidCredentials) {
			slog.Error("LDAP authentication error, falling back to local", "error", err)
		}
	}

	return s.loginLocal(ctx, req)
}

func (s *Service) loginLocal(ctx context.Context, req UserLoginRequest) (*account.Account, error) {
	acc, err := s.accountService.GetAccountByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.ValidatePassword(acc.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.UpdateLastSeen(ctx, acc.ID)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Service) loginLDAP(ctx context.Context, req UserLoginRequest) (*account.Account, error) {
	ldapUser, err := s.ldapService.Authenticate(req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	if s.ldapService.AutoProvision() {
		acc, provisionErr := s.accountService.GetOrCreateLDAPAccount(ctx, ldapUser.Username, ldapUser.Email, ldapUser.Name)
		if provisionErr != nil {
			return nil, provisionErr
		}

		_ = s.accountService.UpdateLastSeen(ctx, acc.ID)
		return acc, nil
	}

	acc, err := s.accountService.GetAccountByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	_ = s.accountService.UpdateLastSeen(ctx, acc.ID)
	return acc, nil
}
