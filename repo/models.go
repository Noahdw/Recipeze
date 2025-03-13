// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthToken struct {
	ID         int32
	Token      string
	UserID     pgtype.Int4
	ConsumedAt pgtype.Timestamptz
	CreatedAt  pgtype.Timestamptz
	ExpiresAt  pgtype.Timestamptz
	CreatorIp  pgtype.Text
}

type Recipe struct {
	ID          int32
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	ImageUrl    pgtype.Text
	Likes       pgtype.Int4
	CreatedAt   pgtype.Timestamptz
}

type User struct {
	ID        int32
	Email     string
	Name      pgtype.Text
	ImageUrl  pgtype.Text
	CreatedAt pgtype.Timestamptz
}
