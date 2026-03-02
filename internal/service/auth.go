package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshData struct {
	Token     string
	UserId    int64
	Id        int64
	ExpiresAt time.Time
}

type UserRepository interface {
	Create(email, login, hashPassword string) error
	GetByCredentials(email, hashPassword string) (int64, error)
}

type TokensRepository interface {
	Get(token string) (RefreshData, error)
	Create(RefreshData) error
}

type PasswordHasher interface {
	Hash(password []byte) (string, error)
}

type AuthService struct {
	secret     string
	hasher     PasswordHasher
	userRepo   UserRepository
	tokensRepo TokensRepository
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

func (s *AuthService) SignIn(email, password string) (string, string, error) {
	hashPass, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return "", "", err
	}

	id, err := s.userRepo.GetByCredentials(email, hashPass)
	if err != nil {
		return "", "", err
	}

	accessToken, err := genAccessToken(id, s.secret, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := genRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = s.tokensRepo.Create(RefreshData{UserId: id, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour * 24 * 30)})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(token string) (string, string, error) {
	tokenData, err := s.tokensRepo.Get(token)
	if err != nil {
		return "", "", err
	}

	if tokenData.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("Token is expires")
	}

	accessToken, err := genAccessToken(tokenData.UserId, s.secret, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := genRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = s.tokensRepo.Create(RefreshData{UserId: tokenData.UserId, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour * 24 * 30)})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
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

func genAccessToken(id int64, secret string, exp time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"Subject": id, "IssuetAt": time.Now().Unix(), "ExpiresAt": time.Now().Add(exp).Unix()})
	return t.SignedString([]byte(secret))
}

func genRefreshToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewAuthService(ur UserRepository, tr TokensRepository, h PasswordHasher, s string) *AuthService {
	return &AuthService{hasher: h, userRepo: ur, secret: s, tokensRepo: tr}
}
