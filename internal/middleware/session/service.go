package middleware

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type MiddlewareService struct {
	repo *MiddlewareRepository
}

func NewMiddlewareService(repository *MiddlewareRepository) *MiddlewareService {
	return &MiddlewareService{repo: repository}
}

func (s *MiddlewareService) FetchSession(sessionID string) (*utils.Session, error) {
	session, err := s.repo.GetSession(sessionID)
	if err != nil || session == nil {
		return nil, errors.New("Session not found")
	}

	expiresAt, err := utils.ConvertTimeStringToRFC3339(session.ExpiresAt)
	if err != nil || time.Now().After(expiresAt) {
		return nil, errors.New("Session expired")
	}

	return session, nil
}

func (s *MiddlewareService) RenewSession(userID string) (*utils.Session, error) {
	newSessionID := uuid.New().String()

	newExpiresAt := time.Now().Add(24 * time.Hour)

	err := s.repo.UpdateSession(userID, newSessionID, newExpiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return &utils.Session{SessionID: newSessionID, ExpiresAt: newExpiresAt.Format(time.RFC3339)}, nil
}
