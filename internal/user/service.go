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
	if !utils.IsEmail(newUser.Email) {
		return fmt.Errorf(`email "%s" is not valid`, newUser.Email)
	}

	if err := s.checkUsernameOrEmailInUse(newUser.Email, newUser.Username); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	newUser.Password = hashedPassword
	newUser.ID = generateID()

	return s.Repo.CreateUser(newUser)
}

func (s *UserService) Login(usernameOrEmail, password string) (string, error) {
	user, err := s.Repo.GetUserByUsernameOrEmail(usernameOrEmail)
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

func (s *UserService) GetUser(userId string) (*UserDetails, error) {

	user, err := s.Repo.GetUserDetailsByID(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserService) checkUsernameOrEmailInUse(email, username string) error {
	usernameExists, err := s.Repo.GetUserByUsername(username)
	if err != nil && err.Error() != "user not found" {
		return fmt.Errorf("error checking username: %v", err)
	}

	emailExists, err := s.Repo.GetUserByEmail(email)
	if err != nil && err.Error() != "user not found" {
		return fmt.Errorf("error checking email: %v", err)
	}

	if usernameExists != nil && emailExists != nil {
		return fmt.Errorf(`username "%s" and email "%s" are taken`, username, email)
	} else if usernameExists != nil {
		return fmt.Errorf(`username "%s" is taken`, username)
	} else if emailExists != nil {
		return fmt.Errorf(`email "%s" is taken`, email)
	}

	return nil
}

func generateID() string {
	id := uuid.New()
	return id.String()
}
