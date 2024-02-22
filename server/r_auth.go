package server

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"libdb.so/ctxt"
)

type authorization struct {
	Token       string
	ParentToken string
	CreatedAt   time.Time
}

func (h *Handler) useAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == nil {
			t, err := h.db.CheckAuthToken(r.Context(), cookie.Value)
			if err != nil {
				h.renderErrorWithCode(w, 404, fmt.Errorf("token not found"))
				return
			}

			token := cookie.Value
			h.logger.Info(
				"Authorizing with token",
				"path", r.URL.Path,
				"token", token,
				"parent", t.ParentToken,
				"created_at", t.CreatedAt)

			ctx := ctxt.With(r.Context(), authorization{
				Token:       token,
				ParentToken: t.ParentToken,
				CreatedAt:   t.CreatedAt,
			})
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := ctxt.From[authorization](r.Context())
		if ok {
			next.ServeHTTP(w, r)
			return
		}

		h.renderErrorWithCode(w, 401, fmt.Errorf("unauthorized"))
	})
}

const dayDuration = 24 * time.Hour

var errTokenNotFound = fmt.Errorf("token not found")

func (h *Handler) getAuth(w http.ResponseWriter, r *http.Request) {
	parentToken := r.PathValue("token")
	newToken := generateNewToken()

	h.logger.Info(
		"Attempting to authenticate for a new token",
		"token", parentToken,
		"new_token", newToken)

	err := h.db.Tx(func(q *sqldb.Queries) error {
		if parentToken == h.rootToken {
			h.logger.Info(
				"Root token authenticated",
				"token", parentToken,
				"new_token", newToken)
		} else {
			_, err := q.CheckAuthToken(r.Context(), parentToken)
			if err != nil {
				return errTokenNotFound
			}

			h.logger.Info(
				"Derived token authenticated",
				"token", parentToken,
				"new_token", newToken)
		}

		_, err := q.AddAuthToken(r.Context(), sqldb.AddAuthTokenParams{
			Token:       newToken,
			ParentToken: parentToken,
		})
		return err
	})
	if err != nil {
		h.logger.Error(
			"Failed to authenticate",
			"token", parentToken,
			"error", err)

		if errors.Is(err, errTokenNotFound) {
			h.renderErrorWithCode(w, 404, err)
		} else {
			h.renderErrorWithCode(w, 500, fmt.Errorf("failed to authenticate"))
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    newToken,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * dayDuration),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateNewToken() string {
	var buf [24]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(buf[:])
}
