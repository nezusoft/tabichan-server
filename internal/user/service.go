package user

import "golang.org/x/crypto/bcrypt"

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Signup(username, password, email string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(username, password string) (*User, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
