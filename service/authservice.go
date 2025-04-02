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
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

type AuthService interface {
	// CreateRegistrationToken generates a secure token for registration
	CreateRegistrationToken(ctx context.Context, email string) (string, error)

	// VerifyRegistrationToken validates a registration token from a request
	VerifyRegistrationToken(ctx context.Context, token string, r *http.Request) (string, error)

	// GetUser retrieves user information by email
	GetUser(ctx context.Context, email string) (*model.User, error)

	// GetUserGroups provides the groups a user belongs to
	GetUserGroups(ctx context.Context, user_id int) ([]model.Group, error)

	// CreateAccount registers a user after verification and sets up default group
	CreateAccount(ctx context.Context, email string) (*model.User, error)
}

func NewAuthService(queries *repo.Queries, db *pgxpool.Pool) *Auth {
	return &Auth{
		queries: queries,
		db:      db,
	}
}

func (a *Auth) CreateRegistrationToken(ctx context.Context, email string) (string, error) {
	token := GenerateSecureToken(32)

	args := repo.CreateRegistrationTokenParams{
		Token: token,
		Email: email,
	}
	err := a.queries.CreateRegistrationToken(ctx, args)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Auth) VerifyRegistrationToken(ctx context.Context, token string, r *http.Request) (string, error) {
	savedToken, err := a.queries.GetRegistrationToken(ctx, token)
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

	args := repo.ConsumeRegistrationTokenParams{
		ConsumedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ID: savedToken.ID,
	}
	err = a.queries.ConsumeRegistrationToken(ctx, args)
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
	pgUser, err := a.queries.GetUserByEmail(ctx, email)
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

func (a *Auth) CreateAccount(ctx context.Context, email string) (*model.User, error) {

	pgUser, err := a.queries.AddUser(ctx, email)
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

func (a *Auth) GetUserGroups(ctx context.Context, user_id int) ([]model.Group, error) {
	pgGroups, err := a.queries.GetUsersGroups(ctx, int32(user_id))
	if err != nil {
		return nil, err
	}
	var groups []model.Group
	for _, g := range pgGroups {
		groups = append(groups, model.Group{
			ID:   int(g.ID),
			Name: g.Name.String,
		})
	}
	return groups, nil
}

func GenerateSecureToken(length int) string {
	// Create a byte slice to store random bytes
	b := make([]byte, length)
	rand.Read(b)

	// Encode to base64 to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)
}
