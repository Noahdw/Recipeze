package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"recipeze/repo"

	"github.com/jackc/pgx/v5/pgtype"
)

type Auth struct {
	db *repo.Queries
}

func (a *Auth) CreateRegistrationToken(ctx context.Context) (string, error) {
	token := GenerateSecureToken(32)

	args := repo.CreateLoginAuthTokenParams{
		Token:  token,
		UserID: pgtype.Int4{},
	}
	err := a.db.CreateLoginAuthToken(ctx, args)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewAuthService(db *repo.Queries) *Auth {
	return &Auth{
		db: db,
	}
}

func GenerateSecureToken(length int) string {
	// Create a byte slice to store random bytes
	b := make([]byte, length)
	rand.Read(b)

	// Encode to base64 to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)
}
