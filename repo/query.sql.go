// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addRecipe = `-- name: AddRecipe :one
INSERT INTO recipes (
    url,
    name,
    description,
    image_url
) VALUES (
    $1, $2, $3, $4
)
RETURNING id
`

type AddRecipeParams struct {
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	ImageUrl    pgtype.Text
}

func (q *Queries) AddRecipe(ctx context.Context, arg AddRecipeParams) (int32, error) {
	row := q.db.QueryRow(ctx, addRecipe,
		arg.Url,
		arg.Name,
		arg.Description,
		arg.ImageUrl,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const addUser = `-- name: AddUser :one
INSERT INTO users (
    email
) VALUES (
    $1
)
RETURNING id
`

func (q *Queries) AddUser(ctx context.Context, email string) (int32, error) {
	row := q.db.QueryRow(ctx, addUser, email)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteRecipeByID = `-- name: DeleteRecipeByID :exec
DELETE FROM recipes where id = $1
`

func (q *Queries) DeleteRecipeByID(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteRecipeByID, id)
	return err
}

const getRecipeByID = `-- name: GetRecipeByID :one
SELECT id, url, name, description, image_url, likes, created_at from recipes WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecipeByID(ctx context.Context, id int32) (Recipe, error) {
	row := q.db.QueryRow(ctx, getRecipeByID, id)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Name,
		&i.Description,
		&i.ImageUrl,
		&i.Likes,
		&i.CreatedAt,
	)
	return i, err
}

const getRecipes = `-- name: GetRecipes :many
SELECT id, url, name, description, image_url, likes, created_at FROM recipes
`

func (q *Queries) GetRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, getRecipes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Name,
			&i.Description,
			&i.ImageUrl,
			&i.Likes,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, name, image_url, created_at from users WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.ImageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const updateRecipe = `-- name: UpdateRecipe :exec
UPDATE recipes 
SET 
    url = $1,
    name = $2,
    description = $3
WHERE id = $4
`

type UpdateRecipeParams struct {
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	ID          int32
}

func (q *Queries) UpdateRecipe(ctx context.Context, arg UpdateRecipeParams) error {
	_, err := q.db.Exec(ctx, updateRecipe,
		arg.Url,
		arg.Name,
		arg.Description,
		arg.ID,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users 
SET 
    image_url = $1,
    name = $2
WHERE id = $3
`

type UpdateUserParams struct {
	ImageUrl pgtype.Text
	Name     pgtype.Text
	ID       int32
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ImageUrl, arg.Name, arg.ID)
	return err
}
