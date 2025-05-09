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
	"recipeze/parsing"
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
		r.Use(m.Authenticate)   // Must be logged in
		r.Use(m.AuthorizeGroup) // Must be member of the group

		// Get recipes for a group
		r.Get("/recipes", h.getRecipes())
		// Add new recipe to a group
		r.Post("/recipes", h.addNewRecipe())
		// Get single recipe (for detail view)
		r.Get("/recipe/{recipe_id}", h.getRecipeDetailView())
		// Show modal for adding a new recipe
		r.Get("/recipes/new", h.showNewRecipeModal())
		// Delete a recipe from a group
		r.Post("/recipes/delete/{recipe_id}", h.deleteRecipe())
		// Get editable details of a recipe
		r.Get("/recipes/update/{recipe_id}", h.updateRecipe())
		// Publish edited details
		r.Post("/recipes/update/{recipe_id}", h.updateRecipeDetails())

		r.Get("/recipes/invite", h.showInviteModal())
		r.Post("/recipes/invite", h.sendInvite())
	})
	// Clear modal
	r.Get("/empty", h.adapt(func(ctx requestContext) (Node, error) { return nil, nil }))
}

func (h *handler) deleteRecipe() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}

		err = h.DeleteRecipeByID(ctx.context(), recipeID)
		if err != nil {
			slog.Error("Could not delete recipe", "ID", recipeID)
			return nil, ErrDefault
		}
		slog.Info("Deleted recipe", "id", recipeID)

		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes")
			return nil, ErrDefault
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

func (h *handler) getRecipes() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			return nil, ErrDefault
		}

		// If HTMX request, return just the list
		if hx.IsRequest(ctx.r.Header) {
			return ui.RecipeListPartial(recipes, 0, groupID), nil
		}
		group := model.Group{
			ID:      groupID,
			Name:    "Test Group",
			Members: []model.GroupMember{},
		}

		users, err := h.GetGroupUsers(ctx.context(), groupID)
		if err != nil {
			return nil, ErrDefault
		}

		for _, user := range users { // Don't like this.
			mg := model.GroupMember{
				ID:      user.ID,
				Name:    user.Name,
				Email:   user.Email,
				IsAdmin: false,
			}
			group.Members = append(group.Members, mg)
		}

		// Otherwise return full page
		return ui.RecipePage(ui.PageProps{IncludeHeader: true}, recipes, group), nil
	})
}

func (h *handler) updateRecipeDetails() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}

		// Parse the form
		err = ctx.r.ParseForm()
		if err != nil {
			return nil, ErrDefault
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
			return nil, ErrDefault
		}

		recipe, err := h.GetRecipeByID(ctx.context(), int32(recipeID))
		if err != nil {
			return nil, ErrDefault
		}

		mainContent := ui.RecipeDetailPartial(recipe, groupID)
		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes")
			return nil, ErrDefault
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
			return nil, ErrDefault
		}
		ctx.r.ParseForm()
		url := ctx.r.FormValue("url")

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}

		resp := req.MustGet(url)

		doc, err := html.Parse(resp.Body)
		if err != nil {
			// ehh, maybe do nothing?
		}
		defer ctx.r.Body.Close()

		meta := extractMeta(doc)
		user := mw.GetUserFromContext(ctx.context())
		id, err := h.AddRecipe(ctx.context(), url, meta.title, meta.description, meta.imageURL, user.ID, groupID)
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, ErrDefault
		}

		//extractedRecipe, err := parsing.ParseRecipe(resp.Bytes())
		//if err != nil {
		//slog.Error("Could not parse recipe", "error", err.Error())

		text := parsing.HtmlToText(resp.Bytes())
		go func() {
			data := parsing.RecipeTextToJsonString(text)
			if data == "" {
				slog.Error("Could not get recipe data from LLM")
				return
			}
			err := h.UpdateRecipeWithJSON(context.Background(), data, id) // FIXME - use better ctx
			if err != nil {
				slog.Error("Could not update db with recipe data", "error", err)
				return
			}
			slog.Info("updated recipe with LLM data", "recipeID", id)
		}()

		//return nil, ErrDefault
		//}
		slog.Info("Added recipe", "userID", user.ID)

		recipe, err := h.GetRecipeByID(ctx.context(), int32(id))
		if err != nil {
			slog.Error("Could not get recipe", "error", err.Error())
			return nil, ErrDefault
		}

		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes", "error", err.Error())
			return nil, ErrDefault
		}

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
			return nil, ErrDefault
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		slog.Info("get recipe", "id", recipeID)

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
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

func (h *handler) updateRecipe() http.HandlerFunc {
	slog.Info("recipe update selected")
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		recipeID, err := getRecipeID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		recipe, err := h.GetRecipeByID(ctx.context(), int32(recipeID))
		if err != nil {
			return nil, ErrDefault
		}

		return ui.RecipeEditPartial(recipe, groupID), nil
	})
}

func (h *handler) showNewRecipeModal() http.HandlerFunc {
	slog.Info("Showed new recipe modal")
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		group_id, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		return ui.RecipeModal(group_id), nil
	})
}

func (h *handler) showInviteModal() http.HandlerFunc {
	slog.Info("Showed invite modal")
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		group_id, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}
		return ui.InviteModal(group_id), nil
	})
}

func (h *handler) sendInvite() http.HandlerFunc {
	slog.Info("Send invite selected")
	return h.adapt(func(ctx requestContext) (Node, error) {
		if !isUserActionAllowed(ctx.context()) {
			return nil, ErrDefault
		}
		ctx.r.ParseForm()

		groupID, err := GetGroupID(ctx.r)
		if err != nil {
			return nil, ErrDefault
		}

		recipes, err := h.GetGroupRecipes(ctx.context(), groupID)
		if err != nil {
			slog.Error("Could not get recipes", "error", err.Error())
			return nil, ErrDefault
		}
		selectedID := 0
		if len(recipes) > 0 {
			selectedID = recipes[0].ID
		}
		mainContent := ui.RecipeListPartial(recipes, selectedID, groupID)
		recipe, err := h.GetRecipeByID(ctx.context(), int32(selectedID))
		if err != nil {
			slog.Error("Could not get recipe", "error", err.Error())
			return nil, ErrDefault
		}
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

func isUserActionAllowed(ctx context.Context) bool {
	authorizedAny := ctx.Value(mw.CtxGroupAuthorizedKey{})
	authorized, ok := authorizedAny.(bool)
	if !ok {
		return false
	}
	return authorized
}
