package service

type UserRepository interface {
	Create(email, login, hashPassword string) error
}

type PasswordHasher interface {
	Hash(password []byte) (string, error)
}

type AuthService struct {
	hasher   PasswordHasher
	userRepo UserRepository
}

func (a *AuthService) SignUp(email, login, password string) error {
	hashPass, err := a.hasher.Hash([]byte(password))
	if err != nil {
		return err
	}

	err = a.userRepo.Create(email, login, hashPass)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) SignIn(email, login, password string) error {
	return nil
}

func NewAuthService(hasher PasswordHasher, userRepo UserRepository) *AuthService {
	return &AuthService{hasher: hasher, userRepo: userRepo}
}
