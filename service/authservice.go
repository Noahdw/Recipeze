package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"recipeze/model"
	"recipeze/repo"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgtype"
)

type Auth struct {
	db *repo.Queries
}

type AuthService interface {
	// CreateRegistrationToken generates a secure token for registration
	CreateRegistrationToken(ctx context.Context, email string) (string, error)

	// VerifyRegistrationToken validates a registration token from a request
	VerifyRegistrationToken(ctx context.Context, token string, r *http.Request) (string, error)

	// GetUser retrieves user information by email
	GetUser(ctx context.Context, email string) (*model.User, error)
}

func NewAuthService(db *repo.Queries) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) CreateRegistrationToken(ctx context.Context, email string) (string, error) {
	token := GenerateSecureToken(32)

	args := repo.CreateLoginAuthTokenParams{
		Token: token,
		Email: email,
	}
	err := a.db.CreateLoginAuthToken(ctx, args)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Auth) VerifyRegistrationToken(ctx context.Context, token string, r *http.Request) (string, error) {
	savedToken, err := a.db.GetLoginAuthToken(ctx, token)
	if err != nil {
		return "", err
	}

	if token != savedToken.Token {
		return "", fmt.Errorf("invalid token")
	}

	if time.Now().After(savedToken.ExpiresAt.Time) {
		return "", fmt.Errorf("expired token")
	}

	if savedToken.ConsumedAt.Valid {
		return "", fmt.Errorf("already consumed")
	}

	args := repo.ConsumeLoginAuthTokenParams{
		ConsumedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ID: savedToken.ID,
	}
	err = a.db.ConsumeLoginAuthToken(ctx, args)
	if err != nil {
		return "", fmt.Errorf("issue consuming token")
	}

	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	print(securecookie.GenerateRandomKey(32))

	session, _ := store.Get(r, "session-name")
	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43

	return savedToken.Email, nil
}

func (a *Auth) GetUser(ctx context.Context, email string) (*model.User, error) {
	pgUser, err := a.db.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}
	user := model.User{
		ID:    int(pgUser.ID),
		Name:  pgUser.Name.String,
		Email: pgUser.Email,
	}
	return &user, nil
}

func GenerateSecureToken(length int) string {
	// Create a byte slice to store random bytes
	b := make([]byte, length)
	rand.Read(b)

	// Encode to base64 to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)
}
