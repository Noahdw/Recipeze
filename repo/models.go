// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Recipe struct {
	ID          int32
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	ImageUrl    pgtype.Text
	CreatedAt   pgtype.Timestamptz
}
