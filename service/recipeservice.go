package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"recipeze/model"
	"recipeze/parsing"
	"recipeze/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/jackc/pgx/v5/pgtype"
)

type Recipe struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewRecipeService(queries *repo.Queries, db *pgxpool.Pool) *Recipe {
	return &Recipe{
		queries: queries,
		db:      db,
	}
}

type RecipeService interface {
	// AddRecipe creates a new recipe and returns its ID
	AddRecipe(ctx context.Context, url, name, description string, imgURL string, userID int, groupID int) (id int, err error)

	// GetRecipes retrieves all recipes
	GetGroupRecipes(ctx context.Context, group_id int) ([]model.Recipe, error)

	// GetRecipeByID retrieves a recipe by its ID
	GetRecipeByID(ctx context.Context, id int32) (*model.Recipe, error)

	// DeleteRecipeByID removes a recipe by its ID
	DeleteRecipeByID(ctx context.Context, id int) error

	// UpdateRecipe modifies an existing recipe
	UpdateRecipe(ctx context.Context, args repo.UpdateRecipeParams) error

	//
	UpdateRecipeWithJSON(ctx context.Context, json string, recipeID int) error
}

func (r *Recipe) AddRecipe(ctx context.Context, url, name, description string, imgURL string, userID int, groupID int) (id int, err error) {
	args := repo.AddRecipeParams{
		CreatedBy:   int32(userID),
		GroupID:     int32(groupID),
		Url:         repo.StringPG(url),
		Name:        repo.StringPG(name),
		Description: repo.StringPG(description),
		ImageUrl:    repo.StringPG(imgURL),
	}
	recipeid, err := r.queries.AddRecipe(ctx, args)
	if err != nil {
		return 0, err
	}

	return int(recipeid), nil
}

func (r *Recipe) UpdateRecipeWithJSON(ctx context.Context, json string, recipeID int) error {
	err := r.queries.UpdateRecipeWithJSON(ctx, repo.UpdateRecipeWithJSONParams{
		DataJson: []byte(json),
		ID:       int32(recipeID),
	})
	return err
}

func (r *Recipe) GetGroupRecipes(ctx context.Context, group_id int) ([]model.Recipe, error) {
	recipesPG, err := r.queries.GetGroupRecipes(ctx, int32(group_id))
	if err != nil {
		return nil, err
	}
	recipes := make([]model.Recipe, 0, len(recipesPG))
	for _, recipePG := range recipesPG {
		recipe := newRecipe(recipePG)
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *Recipe) GetRecipeByID(ctx context.Context, id int32) (*model.Recipe, error) {
	recipePG, err := r.queries.GetRecipeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	recipe := newRecipe(recipePG)
	return &recipe, nil
}

func (r *Recipe) DeleteRecipeByID(ctx context.Context, id int) error {
	err := r.queries.DeleteRecipeByID(ctx, int32(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *Recipe) UpdateRecipe(ctx context.Context, args repo.UpdateRecipeParams) error {
	if len(args.Name.String) == 0 {
		args.Name.String = "Recipe"
	}
	err := r.queries.UpdateRecipe(ctx, args)
	return err
}

func newRecipe(pg repo.Recipe) model.Recipe {
	// Parse the generated JSON
	var collection parsing.RecipeCollection
	err := json.Unmarshal([]byte(pg.DataJson), &collection)
	if err != nil {
		slog.Error("Error unmarshaling recipe json", "error", err)
		file, err := os.Create("recipe.temp")
		if err != nil {
			return model.Recipe{}
		}
		file.Write([]byte(pg.DataJson))
	}

	return model.Recipe{
		ID:          int(pg.ID),
		Name:        pg.Name.String,
		Url:         pg.Url.String,
		Description: pg.Description.String,
		ImageURL:    pg.ImageUrl.String,
		Data:        &collection,
	}
}
