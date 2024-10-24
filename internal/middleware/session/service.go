package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type MiddlewareService struct {
	Repo *MiddlewareRepository
}

var sessionRenewalThreshold = time.Minute * 30

func (s *MiddlewareService) SessionMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		session, err := s.FetchSession(sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		expiresAt, err := utils.ConvertTimeStringToRFC3339(session.ExpiresAt)
		if err != nil {
			http.Error(w, "Session expiry error", http.StatusInternalServerError)
			return
		}

		if time.Until(expiresAt) <= sessionRenewalThreshold {
			newSession, err := s.RenewSession(sessionID)
			if err != nil {
				http.Error(w, "Failed to renew user session", http.StatusInternalServerError)
				return
			}

			newExpiresAt, err := utils.ConvertTimeStringToRFC3339(session.ExpiresAt)
			if err != nil {
				http.Error(w, "Session renewal expiry error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "sessionID",
				Value:    newSession.SessionID,
				Expires:  newExpiresAt,
				HttpOnly: true,
				// Secure: true, // TODO: use Secure when hosting HTTPS
			})
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next(w, r.WithContext(ctx))
	}
}

func (s *MiddlewareService) FetchSession(sessionID string) (*utils.Session, error) {
	session, err := s.Repo.GetSession(sessionID)
	if err != nil || session == nil {
		return nil, errors.New("session not found")
	}

	expiresAt, err := utils.ConvertTimeStringToRFC3339(session.ExpiresAt)
	if err != nil || time.Now().After(expiresAt) {
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (s *MiddlewareService) RenewSession(oldSessionID string) (*utils.Session, error) {
	newSessionID := uuid.New().String()

	newExpiresAt := time.Now().Add(24 * time.Hour)

	err := s.Repo.UpdateSession(oldSessionID, newSessionID, newExpiresAt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return &utils.Session{SessionID: newSessionID, ExpiresAt: newExpiresAt.Format(time.RFC3339)}, nil
}
