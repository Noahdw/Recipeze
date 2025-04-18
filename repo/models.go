// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Group struct {
	ID        int32
	Name      pgtype.Text
	CreatedAt pgtype.Timestamptz
}

type GroupUser struct {
	ID      int32
	GroupID int32
	UserID  int32
}

type LoginToken struct {
	ID         int32
	UserID     int32
	Token      string
	ConsumedAt pgtype.Timestamptz
	CreatedAt  pgtype.Timestamptz
	ExpiresAt  pgtype.Timestamptz
	CreatorIp  pgtype.Text
}

type Recipe struct {
	ID          int32
	CreatedBy   int32
	GroupID     int32
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	DataJson    []byte
	ImageUrl    pgtype.Text
	Likes       pgtype.Int4
	CreatedAt   pgtype.Timestamptz
}

type RegistrationToken struct {
	ID         int32
	Token      string
	Email      string
	ConsumedAt pgtype.Timestamptz
	CreatedAt  pgtype.Timestamptz
	ExpiresAt  pgtype.Timestamptz
	CreatorIp  pgtype.Text
}

type User struct {
	ID           int32
	Email        string
	Name         pgtype.Text
	ImageUrl     pgtype.Text
	SetupAccount pgtype.Bool
	CreatedAt    pgtype.Timestamptz
}
