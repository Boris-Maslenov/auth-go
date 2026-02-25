package service

import (
	"errors"
	"fmt"
	"strconv"
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

func (s *AuthService) ParseToken(token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unixpected signed method %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return 0, errors.New("invalid token")
	}
	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("claims error")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return 0, errors.New("subject error")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func NewAuthService(hasher PasswordHasher, userRepo UserRepository, secret string) *AuthService {
	return &AuthService{hasher: hasher, userRepo: userRepo, secret: secret}
}
