-- name: AddRecipe :one
INSERT INTO recipes (
    url,
    name,
    description
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: GetRecipes :many
SELECT * FROM recipes;

-- name: GetRecipeByID :one
SELECT * from recipes WHERE id = $1 LIMIT 1;

-- name: DeleteRecipeByID :exec
DELETE FROM recipes where id = $1;