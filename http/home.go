package http

import (
	//"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	hx "maragu.dev/gomponents-htmx/http"
	ghttp "maragu.dev/gomponents/http"

	"recipeze/html"
	"recipeze/service"
)

func Home(r chi.Router, s *service.Recipe) {
	// Get single recipe (for detail view)
	r.Get("/recipes/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")
		slog.Info("get recipe", "id", idstr)
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}

		recipe, err := s.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return html.ErrorPartial("Recipe not found"), nil
		}

		return html.RecipeDetailPartial(recipe), nil
	}))

	// Get all recipes (main page)
	r.Get("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// If HTMX request, return just the list
		if hx.IsRequest(r.Header) {
			return html.RecipeListPartial(recipes, time.Now()), nil
		}

		// Otherwise return full page
		return html.RecipePage(html.PageProps{}, recipes, time.Now()), nil
	}))

	// Add new recipe
	r.Post("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		err := s.AddRecipe(r.Context(), "https://example.com", "New Recipe", "A tasty dish")
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// Return just the updated list
		return html.RecipeListPartial(recipes, time.Now()), nil
	}))
}
