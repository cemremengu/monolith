package auth

import (
	"context"

	"monolith/internal/database"
	"monolith/internal/service/account"
)

type Service struct {
	db             *database.DB
	accountService *account.Service
}

func NewService(db *database.DB) *Service {
	return &Service{
		db:             db,
		accountService: account.NewService(db),
	}
}

func (s *Service) Register(ctx context.Context, req account.RegisterRequest) (*account.Account, error) {
	if req.Password == "" || len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	exists, err := s.accountService.UserExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	return s.accountService.Register(ctx, req)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*account.Account, error) {
	account, err := s.accountService.GetAccountByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.ValidatePassword(account.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = s.accountService.UpdateLastSeen(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
