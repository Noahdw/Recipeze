package service

import (
	"context"
	"recipeze/model"
	"recipeze/repository"
)

type Recipe struct {
	db *repository.Queries
}

func NewRecipeService(db *repository.Queries) *Recipe {
	return &Recipe{
		db: db,
	}
}

func (r *Recipe) AddRecipe(ctx context.Context, url, name, description string) error {
	args := repository.AddRecipeParams{
		Url:         repository.StringPG(url),
		Name:        repository.StringPG(name),
		Description: repository.StringPG(description),
	}
	_, err := r.db.AddRecipe(ctx, args)
	if err != nil {
		return err
	}
	return nil
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
		Name: recipePG.Name.String,
		Url:  recipePG.Url.String,
	}

	return recipe, nil
}
