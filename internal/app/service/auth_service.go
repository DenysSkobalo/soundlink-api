package service

import (
	"errors"

	"github.com/DenysSkobalo/soundlink-api/internal/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Signup(email, password string) (map[string]string, error) {
	if email == "" || password == "" {
		return nil, errors.New("All fields are required")
	}

	if !utils.IsValidEmail(email) {
		return nil, errors.New("Invalid email format")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}

	return map[string]string{
		"email":    email,
		"password": hashedPassword,
	}, nil
}
