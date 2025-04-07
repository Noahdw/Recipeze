package handler

import (
	//"context"

	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	mw "recipeze/middleware"
	"recipeze/model"
	"recipeze/repo"
	"recipeze/ui"

	"github.com/imroc/req/v3"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	hx "maragu.dev/gomponents-htmx/http"
	. "maragu.dev/gomponents/html"
)

type meta struct {
	title       string
	description string
	siteName    string
	imageURL    string
}

func (h *handler) RouteRecipe(r chi.Router, m *mw.AuthMiddleware) {
	r.Route("/g/{group_id}", func(r chi.Router) {
		r.Use(m.Authenticate)
		r.Use(m.AuthorizeGroup)

		// Get page for a group, including recipes
		r.Get("/recipes", h.adapt(func(ctx requestContext) (Node, error) {
			if !isUserActionAllowed(ctx.context()) {
				return nil, fmt.Errorf("")
			}
			groupID, err := GetGroupID(ctx.r)
			if err != nil {
				return nil, err
			}
			recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
			if err != nil {
				return nil, err
			}

			// If HTMX request, return just the list
			if hx.IsRequest(ctx.r.Header) {
				return ui.RecipeListPartial(recipes, 0, groupID), nil
			}

			// Otherwise return full page
			return ui.RecipePage(ui.PageProps{IncludeHeader: true}, recipes, groupID), nil
		}))
		// Add new recipe
		r.Post("/recipes", h.addNewRecipe())

		// Get single recipe (for detail view)
		r.Get("/recipe/{recipe_id}", h.getRecipeDetailView())

		r.Get("/recipes/new", h.adapt(func(ctx requestContext) (Node, error) {
			if !isUserActionAllowed(ctx.context()) {
				return nil, fmt.Errorf("")
			}
			group_id, err := GetGroupID(ctx.r)
			if err != nil {
				return nil, err
			}
			return ui.RecipeModal(group_id), nil
		}))

		r.Post("/recipes/delete/{recipe_id}", h.deleteRecipe())

		r.Get("/recipes/update/{recipe_id}", h.adapt(func(ctx requestContext) (Node, error) {
			if !isUserActionAllowed(ctx.context()) {
				return nil, fmt.Errorf("")
			}
			recipeID, err := getRecipeID(ctx.r)
			if err != nil {
				return nil, err
			}
			groupID, err := GetGroupID(ctx.r)
			if err != nil {
				return nil, err
			}
			recipe, err := h.GetRecipeByID(ctx.context(), int32(recipeID))
			if err != nil {
				return nil, err
			}

			return ui.RecipeEditPartial(recipe, groupID), nil
		}))

		r.Post("/recipes/update/{recipe_id}", h.updateRecipeDetails())
	})

	r.Get("/empty", h.adapt(func(ctx requestContext) (Node, error) {
		// Clear modal
		return nil, nil
	}))
}

func (h *handler) deleteRecipe() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, fmt.Errorf("")
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, err
		}
		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, err
		}

		err = h.DeleteRecipeByID(ctx.context(), recipeID)
		if err != nil {
			slog.Error("Could not delete recipe", "ID", recipeID)
			return nil, err
		}
		slog.Info("Deleted recipe", "id", recipeID)

		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
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

		mainContent := ui.RecipeDetailPartial(&recipe, groupID)

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-list"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListPartial(recipes, selectedID, groupID),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) updateRecipeDetails() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, fmt.Errorf("")
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, err
		}

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, err
		}

		// Parse the form
		err = ctx.r.ParseForm()
		if err != nil {
			return nil, err
		}

		// Update the recipe in the database
		err = h.UpdateRecipe(ctx.context(), repo.UpdateRecipeParams{
			ID:          int32(recipeID),
			Name:        repo.StringPG(ctx.r.FormValue("name")),
			Url:         repo.StringPG(ctx.r.FormValue("url")),
			Description: repo.StringPG(ctx.r.FormValue("description")),
		})
		if err != nil {
			slog.Error("Could not update recipe", "ID", recipeID)
			return nil, err
		}

		recipe, err := h.GetRecipeByID(ctx.context(), int32(recipeID))
		if err != nil {
			return nil, err
		}

		mainContent := ui.RecipeDetailPartial(recipe, groupID)
		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes")
			return nil, err
		}

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-list"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListPartial(recipes, recipeID, groupID),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) addNewRecipe() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, fmt.Errorf("")
		}
		ctx.r.ParseForm()
		url := ctx.r.FormValue("url")

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, err
		}

		resp := req.MustGet(url)

		doc, err := html.Parse(resp.Body)
		if err != nil {
			// ehh, maybe do nothing?
		}
		defer ctx.r.Body.Close()

		meta := extractMeta(doc)
		id, err := h.AddRecipe(ctx.context(), url, meta.title, meta.description, meta.imageURL, 1, groupID) //FIX
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipe, err := h.GetRecipeByID(ctx.context(), int32(id))

		if err != nil {
			slog.Error("Could not get recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes", "error", err.Error())
			return nil, err
		}
		fmt.Printf("%#v\n", recipes)
		mainContent := ui.RecipeListPartial(recipes, id, groupID)

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-detail"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeDetailPartial(recipe, groupID),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func (h *handler) getRecipeDetailView() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, fmt.Errorf("")
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, err
		}
		slog.Info("get recipe", "id", recipeID)

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, err
		}

		recipe, err := h.GetRecipeByID(ctx.r.Context(), int32(recipeID))
		if err != nil {
			return ui.ErrorPartial("Recipe not found"), nil
		}

		mainContent := ui.RecipeDetailPartial(recipe, groupID)

		//recipes, err := s.GetRecipes(r.Context())
		listItemID := fmt.Sprintf("recipe-list-item-%d", recipe.ID)
		// Second part updates another element out-of-band
		listContent := Div(
			ID(listItemID),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListItemPartial(recipe, recipeID, groupID),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil
	})
}

func isUserActionAllowed(ctx context.Context) bool {
	authorizedAny := ctx.Value(mw.CtxGroupAuthorizedKey{})
	authorized, ok := authorizedAny.(bool)
	if !ok {
		return false
	}
	return authorized
}
