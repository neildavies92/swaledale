package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/neildavies/swaledale/internal/auth"
	"github.com/neildavies/swaledale/internal/domain"
	"github.com/neildavies/swaledale/internal/store"
)

const (
	sessionCookie   = "swaledale_session"
	sessionDuration = 30 * 24 * time.Hour
)

type contextKey string

const userContextKey contextKey = "user"

type registerRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	HouseholdName string `json:"householdName"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) register(w http.ResponseWriter, r *http.Request) {
	var input registerRequest
	if !decode(w, r, &input) {
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.HouseholdName = strings.TrimSpace(input.HouseholdName)
	switch {
	case input.Name == "":
		writeError(w, http.StatusBadRequest, "name is required")
		return
	case !strings.Contains(input.Email, "@"):
		writeError(w, http.StatusBadRequest, "a valid email is required")
		return
	case len(input.Password) < 8:
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if input.HouseholdName == "" {
		input.HouseholdName = input.Name + "'s Household"
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	user, err := api.store.RegisterUser(r.Context(), store.RegisterUserInput{
		Name:          input.Name,
		Email:         input.Email,
		PasswordHash:  passwordHash,
		HouseholdName: input.HouseholdName,
	})
	if errors.Is(err, store.ErrEmailTaken) {
		writeError(w, http.StatusConflict, "that email is already registered")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.startSession(w, r, user)
}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest
	if !decode(w, r, &input) {
		return
	}
	user, passwordHash, err := api.store.UserByEmail(r.Context(), strings.TrimSpace(input.Email))
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !auth.CheckPassword(passwordHash, input.Password)) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.startSession(w, r, user)
}

func (api *API) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if err := api.store.DeleteSession(r.Context(), auth.HashSessionToken(cookie.Value)); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	http.SetCookie(w, api.sessionCookie("", -1))
	writeJSON(w, http.StatusOK, map[string]string{"status": "signed out"})
}

func (api *API) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, userFrom(r.Context()))
}

func (api *API) startSession(w http.ResponseWriter, r *http.Request, user domain.User) {
	token, err := auth.NewSessionToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	if err := api.store.CreateSession(r.Context(), auth.HashSessionToken(token), user.ID, time.Now().Add(sessionDuration)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(w, api.sessionCookie(token, int(sessionDuration.Seconds())))
	writeJSON(w, http.StatusOK, user)
}

func (api *API) sessionCookie(token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   api.secureCookies,
		SameSite: http.SameSiteLaxMode,
	}
}

func (api *API) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil || cookie.Value == "" {
			writeError(w, http.StatusUnauthorized, "sign in required")
			return
		}
		user, err := api.store.UserBySession(r.Context(), auth.HashSessionToken(cookie.Value))
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "session expired")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func userFrom(ctx context.Context) domain.User {
	user, _ := ctx.Value(userContextKey).(domain.User)
	return user
}
