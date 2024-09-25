package user

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type UserService struct {
	Repo *UserRepository
}

func (s *UserService) Signup(newUser User) error {
	exists, err := s.Repo.GetUserByUsername(newUser.Username)
	if err != nil && err.Error() != "user not found" {
		return fmt.Errorf("error checking username: %v", err)
	}

	if exists != nil {
		return fmt.Errorf(`username "%s" is taken`, newUser.Username)
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	newUser.Password = hashedPassword
	newUser.ID = generateID()

	return s.Repo.CreateUser(newUser)
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
