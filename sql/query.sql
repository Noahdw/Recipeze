-- name: AddRecipe :one
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
RETURNING id;

-- name: GetGroupRecipes :many 
SELECT r.* 
FROM recipes r
JOIN group_recipes gr ON r.id = gr.recipe_id
WHERE gr.group_id = $1;

-- name: GetUserRecipes :many
SELECT * FROM recipes where created_by = $1;

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

-- name: GetGroupUsers :many
SELECT u.* 
FROM users u
JOIN group_users gu ON u.id = gu.user_id
WHERE gu.group_id = $1;

-- name: GetUserByEmail :one
SELECT * from users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users 
SET 
    image_url = $1,
    name = $2
WHERE id = $3;

-- name: CreateRegistrationToken :exec
INSERT INTO registration_tokens (
    token,
    email
) VALUES (
    $1, $2
);

-- name: GetRegistrationToken :one
SELECT * FROM registration_tokens WHERE token = $1 LIMIT 1;

-- name: ConsumeRegistrationToken :exec
UPDATE registration_tokens
SET
    consumed_at = $1
WHERE id = $2;