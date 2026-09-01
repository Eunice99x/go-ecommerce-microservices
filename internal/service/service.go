package service

import "github.com/eunice99x/goMicro/internal/pkg/auth"

type Service struct {
	storer   Storer
	tokenGen *auth.JWTConfig
}

func NewService(storer Storer, tokenGen *auth.JWTConfig) *Service {
	return &Service{
		storer:   storer,
		tokenGen: tokenGen,
	}
}
