-- name: AddRecipe :execresult
INSERT INTO recipes (
    url,
    name,
    description
) VALUES (
    $1, $2, $3
);

-- name:   :many
SELECT * from recipes;

-- name: GetRecipeByID :one
select * from recipes where id = $1 LIMIT 1;