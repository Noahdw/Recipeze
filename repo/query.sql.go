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
    created_by,
    group_id,
    url,
    name,
    description,
    image_url
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id
`

type AddRecipeParams struct {
	CreatedBy   int32
	GroupID     int32
	Url         pgtype.Text
	Name        pgtype.Text
	Description pgtype.Text
	ImageUrl    pgtype.Text
}

func (q *Queries) AddRecipe(ctx context.Context, arg AddRecipeParams) (int32, error) {
	row := q.db.QueryRow(ctx, addRecipe,
		arg.CreatedBy,
		arg.GroupID,
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
RETURNING id, email, name, image_url, created_at
`

func (q *Queries) AddUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, addUser, email)
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

const addUserToGroup = `-- name: AddUserToGroup :exec
INSERT INTO group_users (
    group_id,
    user_id
) VALUES (
    $1, $2
)
`

type AddUserToGroupParams struct {
	GroupID int32
	UserID  int32
}

func (q *Queries) AddUserToGroup(ctx context.Context, arg AddUserToGroupParams) error {
	_, err := q.db.Exec(ctx, addUserToGroup, arg.GroupID, arg.UserID)
	return err
}

const consumeRegistrationToken = `-- name: ConsumeRegistrationToken :exec
UPDATE registration_tokens
SET
    consumed_at = $1
WHERE id = $2
`

type ConsumeRegistrationTokenParams struct {
	ConsumedAt pgtype.Timestamptz
	ID         int32
}

func (q *Queries) ConsumeRegistrationToken(ctx context.Context, arg ConsumeRegistrationTokenParams) error {
	_, err := q.db.Exec(ctx, consumeRegistrationToken, arg.ConsumedAt, arg.ID)
	return err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO groups (
    name
) VALUES (
    $1
)
RETURNING id, name, created_at
`

func (q *Queries) CreateGroup(ctx context.Context, name pgtype.Text) (Group, error) {
	row := q.db.QueryRow(ctx, createGroup, name)
	var i Group
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const createRegistrationToken = `-- name: CreateRegistrationToken :exec
INSERT INTO registration_tokens (
    token,
    email
) VALUES (
    $1, $2
)
`

type CreateRegistrationTokenParams struct {
	Token string
	Email string
}

func (q *Queries) CreateRegistrationToken(ctx context.Context, arg CreateRegistrationTokenParams) error {
	_, err := q.db.Exec(ctx, createRegistrationToken, arg.Token, arg.Email)
	return err
}

const deleteRecipeByID = `-- name: DeleteRecipeByID :exec
DELETE FROM recipes where id = $1
`

func (q *Queries) DeleteRecipeByID(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteRecipeByID, id)
	return err
}

const getGroupRecipes = `-- name: GetGroupRecipes :many
SELECT id, created_by, group_id, url, name, description, image_url, likes, created_at FROM recipes where group_id = $1
`

func (q *Queries) GetGroupRecipes(ctx context.Context, groupID int32) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, getGroupRecipes, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.CreatedBy,
			&i.GroupID,
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

const getGroupUsers = `-- name: GetGroupUsers :many
SELECT u.id, u.email, u.name, u.image_url, u.created_at 
FROM users u
JOIN group_users gu ON u.id = gu.user_id
WHERE gu.group_id = $1
`

func (q *Queries) GetGroupUsers(ctx context.Context, groupID int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getGroupUsers, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Name,
			&i.ImageUrl,
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

const getRecipeByID = `-- name: GetRecipeByID :one
SELECT id, created_by, group_id, url, name, description, image_url, likes, created_at from recipes WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecipeByID(ctx context.Context, id int32) (Recipe, error) {
	row := q.db.QueryRow(ctx, getRecipeByID, id)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.CreatedBy,
		&i.GroupID,
		&i.Url,
		&i.Name,
		&i.Description,
		&i.ImageUrl,
		&i.Likes,
		&i.CreatedAt,
	)
	return i, err
}

const getRegistrationToken = `-- name: GetRegistrationToken :one
SELECT id, token, email, consumed_at, created_at, expires_at, creator_ip FROM registration_tokens WHERE token = $1 LIMIT 1
`

func (q *Queries) GetRegistrationToken(ctx context.Context, token string) (RegistrationToken, error) {
	row := q.db.QueryRow(ctx, getRegistrationToken, token)
	var i RegistrationToken
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.Email,
		&i.ConsumedAt,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.CreatorIp,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, name, image_url, created_at from users WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
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

const getUserRecipes = `-- name: GetUserRecipes :many
SELECT id, created_by, group_id, url, name, description, image_url, likes, created_at FROM recipes where created_by = $1
`

func (q *Queries) GetUserRecipes(ctx context.Context, createdBy int32) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, getUserRecipes, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.CreatedBy,
			&i.GroupID,
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

const getUsersGroups = `-- name: GetUsersGroups :many
SELECT groups.id, groups.name, groups.created_at 
FROM groups
JOIN group_users ON groups.id = group_users.user_id
WHERE group_users.group_id = $1
`

func (q *Queries) GetUsersGroups(ctx context.Context, groupID int32) ([]Group, error) {
	rows, err := q.db.Query(ctx, getUsersGroups, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Group
	for rows.Next() {
		var i Group
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
