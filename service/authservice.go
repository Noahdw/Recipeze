package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"recipeze/model"
	"recipeze/repo"
	"time"

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

	// Login logs in a user and returns a session token
	Login(ctx context.Context, userID int, token string) (bool, error)

	// GetLoggedInUser gives the user from a session token if they are logged in
	GetLoggedInUser(ctx context.Context, session_token string) (*model.User, error)

	// IsUserInGroup tells if a user belongs to a group
	IsUserInGroup(ctx context.Context, groupID int, userID int) (bool, error)
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
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := a.queries.WithTx(tx)

	pgUser, err := qtx.AddUser(ctx, email)
	if err != nil {
		return nil, err
	}

	pgGroup, err := qtx.CreateGroup(ctx, repo.StringPG("Your recipes"))
	if err != nil {
		return nil, err
	}

	err = qtx.AddUserToGroup(ctx, repo.AddUserToGroupParams{
		GroupID: pgGroup.ID,
		UserID:  pgUser.ID,
	})
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:    int(pgUser.ID),
		Name:  pgUser.Name.String,
		Email: pgUser.Email,
	}
	return &user, tx.Commit(ctx)
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

func (a *Auth) Login(ctx context.Context, userID int, token string) (bool, error) {
	_, err := a.queries.CreateLoginToken(ctx, repo.CreateLoginTokenParams{
		UserID: int32(userID),
		Token:  token,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *Auth) GetLoggedInUser(ctx context.Context, sessionStoken string) (*model.User, error) {
	token, err := a.queries.GetLoginToken(ctx, sessionStoken)
	if err != nil {
		return nil, err
	}

	if time.Now().After(token.ExpiresAt.Time) {
		return nil, fmt.Errorf("")
	}
	pgUser, err := a.queries.GetUserByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:    int(pgUser.ID),
		Name:  pgUser.Name.String,
		Email: pgUser.Email,
	}
	return user, nil
}

func (a *Auth) IsUserInGroup(ctx context.Context, groupID int, userID int) (bool, error) {
	_, err := a.queries.IsUserInGroup(ctx, repo.IsUserInGroupParams{
		GroupID: int32(groupID),
		UserID:  int32(userID),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateSecureToken(length int) string {
	// Create a byte slice to store random bytes
	b := make([]byte, length)
	rand.Read(b)

	// Encode to base64 to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)
}
