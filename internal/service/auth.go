package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRepository interface {
	Create(email, login, hashPassword string) error
	GetByCredentials(email, hashPassword string) (string, error)
}

type PasswordHasher interface {
	Hash(password []byte) (string, error)
}

type AuthService struct {
	secret   string
	hasher   PasswordHasher
	userRepo UserRepository
}

func (s *AuthService) SignUp(email, login, password string) error {
	hashPass, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return err
	}

	err = s.userRepo.Create(email, login, hashPass)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(email, password string) (string, error) {
	hashPass, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return "", err
	}

	id, err := s.userRepo.GetByCredentials(email, hashPass)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"Subject": id, "IssuetAt": time.Now().Unix(), "ExpiresAt": time.Now().Add(time.Minute * 15).Unix()})
	return token.SignedString([]byte(s.secret))
}

func NewAuthService(hasher PasswordHasher, userRepo UserRepository, secret string) *AuthService {
	return &AuthService{hasher: hasher, userRepo: userRepo, secret: secret}
}
