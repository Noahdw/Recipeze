-- name: AddRecipe :one
INSERT INTO recipes (
    url,
    name,
    description,
    image_url
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: GetRecipes :many
SELECT * FROM recipes;

-- name: GetRecipeByID :one
SELECT * from recipes WHERE id = $1 LIMIT 1;

-- name: DeleteRecipeByID :exec
DELETE FROM recipes where id = $1;

-- name: UpdateRecipe :exec
UPDATE recipes 
SET 
    url = $1,
    name = $2,
    description = $3
WHERE id = $4;

-- name: AddUser :one
INSERT INTO users (
    email
) VALUES (
    $1
)
RETURNING id;

-- name: GetUser :one
SELECT * from users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users 
SET 
    image_url = $1,
    name = $2
WHERE id = $3;

-- name: CreateLoginAuthToken :exec
INSERT INTO auth_tokens (
    token,
    email
) VALUES (
    $1, $2
);

-- name: GetLoginAuthToken :one
SELECT * FROM auth_tokens WHERE token = $1 LIMIT 1;

-- name: ConsumeLoginAuthToken :exec
UPDATE auth_tokens
SET
    consumed_at = $1
WHERE id = $2;