package login

import (
	"context"

	"monolith/internal/database"
	"monolith/internal/service/account"
)

type Service struct {
	db             *database.DB
	accountService *account.Service
}

func NewService(db *database.DB, accountService *account.Service) *Service {
	return &Service{
		db:             db,
		accountService: accountService,
	}
}

func (s *Service) Login(ctx context.Context, req UserLoginRequest) (*account.Account, error) {
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
