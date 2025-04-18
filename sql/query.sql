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

-- name: UpdateRecipeWithJSON :exec
UPDATE recipes 
SET 
    data_json = $1
WHERE id = $2;

-- name: GetGroupRecipes :many 
SELECT * FROM recipes where group_id = $1;

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
RETURNING *;

-- name: CreateGroup :one
INSERT INTO groups (
    name
) VALUES (
    $1
)
RETURNING *;

-- name: AddUserToGroup :exec
INSERT INTO group_users (
    group_id,
    user_id
) VALUES (
    $1, $2
);

-- name: GetGroupUsers :many
SELECT u.* 
FROM users u
JOIN group_users gu ON u.id = gu.user_id
WHERE gu.group_id = $1;

-- name: GetUsersGroups :many
SELECT groups.* 
FROM groups
JOIN group_users ON groups.id = group_users.user_id
WHERE group_users.group_id = $1;

-- name: IsUserInGroup :one
SELECT id from group_users WHERE group_id = $1 AND user_id = $2 LIMIT 1;

-- name: IsUserAccountSetupComplete :one
select setup_account from users where id = $1;

-- name: GetUserByEmail :one
SELECT * from users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * from users WHERE id = $1 LIMIT 1;

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

-- name: CreateLoginToken :one
INSERT INTO login_tokens (
    user_id,
    token,
    expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetLoginToken :one
SELECT * FROM login_tokens WHERE token = $1 LIMIT 1;