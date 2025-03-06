package service

import (
	"context"
	"recipeze/model"
	"recipeze/repo"
)

type Recipe struct {
	db *repo.Queries
}

func NewRecipeService(db *repo.Queries) *Recipe {
	return &Recipe{
		db: db,
	}
}

func (r *Recipe) AddRecipe(ctx context.Context, url, name, description string) (id int, err error) {
	args := repo.AddRecipeParams{
		Url:         repo.StringPG(url),
		Name:        repo.StringPG(name),
		Description: repo.StringPG(description),
	}
	recipeid, err := r.db.AddRecipe(ctx, args)
	if err != nil {
		return 0, err
	}

	return int(recipeid), nil
}

func (r *Recipe) GetRecipes(ctx context.Context) ([]model.Recipe, error) {
	recipesPG, err := r.db.GetRecipes(ctx)
	if err != nil {
		return nil, err
	}
	recipes := make([]model.Recipe, 0, len(recipesPG))
	for _, recipePG := range recipesPG {
		recipe := model.Recipe{
			Name:        recipePG.Name.String,
			Url:         recipePG.Url.String,
			ID:          int(recipePG.ID),
			Description: recipePG.Description.String,
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *Recipe) GetRecipeByID(ctx context.Context, id int32) (*model.Recipe, error) {
	recipePG, err := r.db.GetRecipeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	recipe := &model.Recipe{
		ID:   int(id),
		Name: recipePG.Name.String,
		Url:  recipePG.Url.String,
	}

	return recipe, nil
}

func (r *Recipe) DeleteRecipeByID(ctx context.Context, id int) error {
	err := r.db.DeleteRecipeByID(ctx, int32(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *Recipe) UpdateRecipe(ctx context.Context, args repo.UpdateRecipeParams) error {
	err := r.db.UpdateRecipe(ctx, args)
	return err
}
