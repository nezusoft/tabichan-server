package middleware

import (
	"net/http"
	"time"

	"github.com/tabichanorg/tabichan-server/internal/utils"
)

var sessionRenewalThreshold = time.Minute * 30

func SessionMiddleware(middlewareService *MiddlewareService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("sessionID")
			if err != nil {
				http.Error(w, "Session not found", http.StatusUnauthorized)
				return
			}

			sessionID := cookie.Value

			session, err := utils.FetchSession(sessionID)
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
				newSession, err := utils.RenewSession(session.UserID)
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
					// Secure:   true, // TODO: use Secure when hosting HTTPS
				})
			}

			next.ServeHTTP(w, r)
		})
	}
}
