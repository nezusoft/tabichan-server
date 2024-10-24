package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type UserService struct {
	Repo *UserRepository
}

func (s *UserService) Signup(newUser UserLogin, device string) (*LoginRequestResponse, error) {
	if !utils.IsEmail(newUser.Email) {
		return &LoginRequestResponse{}, fmt.Errorf(`email "%s" is not valid`, newUser.Email)
	}

	if err := s.checkUsernameOrEmailInUse(newUser.Email, newUser.Username); err != nil {
		return &LoginRequestResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return &LoginRequestResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	newUser.Password = hashedPassword
	newUser.UserID = generateID()

	expiresAt := time.Now().Add(24 * time.Hour)

	session, err := s.createNewSession(expiresAt, newUser.UserID, device)
	if err != nil {
		return &LoginRequestResponse{}, err
	}

	token, err := utils.GenerateJWT(newUser.Username)
	if err != nil {
		return &LoginRequestResponse{}, fmt.Errorf("failed to generate token: %v", err)
	}

	return &LoginRequestResponse{Token: token, Session: session}, s.Repo.CreateUser(newUser)
}

func (s *UserService) Login(usernameOrEmail, password, device string, rememberMeSelected bool) (*LoginRequestResponse, error) {
	user, err := s.Repo.GetUserByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		return &LoginRequestResponse{}, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return &LoginRequestResponse{}, fmt.Errorf("invalid username or password")
	}

	var expiresAt time.Time
	if rememberMeSelected {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	} else {
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	session, err := s.createNewSession(expiresAt, user.UserID, device)
	if err != nil {
		return &LoginRequestResponse{}, err
	}

	// TODO: incorporate before MVP completion
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return &LoginRequestResponse{}, fmt.Errorf("failed to generate token: %v", err)
	}

	return &LoginRequestResponse{Token: token, Session: session}, nil
}

func (s *UserService) GetUser(userId string) (*User, error) {

	user, err := s.Repo.GetUserDetailsByID(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserService) createNewSession(expiresAt time.Time, userID string, device string) (*utils.Session, error) {
	session := &utils.Session{
		SessionID: uuid.New().String(),
		UserID:    userID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		CreatedAt: time.Now().Format(time.RFC3339),
		Device:    device,
	}

	err := s.Repo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
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
