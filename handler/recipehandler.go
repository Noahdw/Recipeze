package handler

import (
	//"context"

	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	"recipeze/model"
	"recipeze/repo"
	"recipeze/ui"

	"github.com/imroc/req/v3"
	"golang.org/x/net/html"
	hx "maragu.dev/gomponents-htmx/http"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type meta struct {
	title       string
	description string
	siteName    string
	imageURL    string
}

func RouteRecipe(r chi.Router, handler *handler) {
	// Get single recipe (for detail view)
	r.Get("/recipes/{id}", handler.getRecipeDetailView())

	// Get all recipes (main page)
	r.Get("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		recipes, err := handler.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// If HTMX request, return just the list
		if hx.IsRequest(r.Header) {
			return ui.RecipeListPartial(recipes, 0), nil
		}

		// Otherwise return full page
		return ui.RecipePage(ui.PageProps{}, recipes), nil
	}))

	// Add new recipe
	r.Post("/recipes", handler.addNewRecipe())

	r.Get("/recipes/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Show modal
		return ui.RecipeModal(), nil
	}))

	r.Post("/recipes/delete/{id}", handler.deleteRecipe())

	r.Get("/recipes/update/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}
		recipe, err := handler.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return nil, err
		}

		return ui.RecipeEditPartial(recipe), nil
	}))

	r.Post("/recipes/update/{id}", handler.updateRecipeDetails())

	r.Get("/empty", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Clear modal
		return nil, nil
	}))
}

func (h *handler) deleteRecipe() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}

		err = h.DeleteRecipeByID(r.Context(), id)
		if err != nil {
			slog.Error("Could not delete recipe", "ID", id)
			return nil, err
		}
		slog.Info("Deleted recipe", "id", id)

		recipes, err := h.GetRecipes(r.Context())
		if err != nil {
			slog.Error("Could not get recipes")
			return nil, err
		}
		var recipe model.Recipe
		selectedID := 0
		if len(recipes) > 0 {
			recipe = recipes[0]
			selectedID = recipe.ID
		}

		mainContent := ui.RecipeDetailPartial(&recipe)

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-list"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListPartial(recipes, selectedID),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) updateRecipeDetails() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}

		// Parse the form
		err = r.ParseForm()
		if err != nil {
			return nil, err
		}

		// Update the recipe in the database
		err = h.UpdateRecipe(r.Context(), repo.UpdateRecipeParams{
			ID:          int32(id),
			Name:        repo.StringPG(r.FormValue("name")),
			Url:         repo.StringPG(r.FormValue("url")),
			Description: repo.StringPG(r.FormValue("description")),
		})
		if err != nil {
			slog.Error("Could not update recipe", "ID", id)
			return nil, err
		}

		recipe, err := h.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return nil, err
		}

		mainContent := ui.RecipeDetailPartial(recipe)
		recipes, err := h.GetRecipes(r.Context())
		if err != nil {
			slog.Error("Could not get recipes")
			return nil, err
		}

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-list"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListPartial(recipes, id),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) addNewRecipe() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		r.ParseForm()
		url := r.FormValue("url")

		resp := req.MustGet(url)

		doc, err := html.Parse(resp.Body)
		if err != nil {
			// ehh, maybe do nothing?
		}
		defer r.Body.Close()

		meta := extractMeta(doc)
		id, err := h.AddRecipe(r.Context(), url, meta.title, meta.description, meta.imageURL)
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := h.GetRecipes(r.Context())
		if err != nil {
			slog.Error("Could not get recipes", "error", err.Error())
			return nil, err
		}

		recipe, err := h.GetRecipeByID(r.Context(), int32(id))

		if err != nil {
			slog.Error("Could not get recipe", "error", err.Error())
			return nil, err
		}
		mainContent := ui.RecipeListPartial(recipes, id)

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-detail"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeDetailPartial(recipe),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) getRecipeDetailView() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")
		slog.Info("get recipe", "id", idstr)
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}

		recipe, err := h.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return ui.ErrorPartial("Recipe not found"), nil
		}

		mainContent := ui.RecipeDetailPartial(recipe)

		//recipes, err := s.GetRecipes(r.Context())
		listItemID := fmt.Sprintf("recipe-list-item-%d", recipe.ID)
		// Second part updates another element out-of-band
		listContent := Div(
			ID(listItemID),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListItemPartial(recipe, id),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil

	})
}
