package http

import (
	//"context"

	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	"recipeze/model"
	"recipeze/service"
	"recipeze/ui"

	"github.com/imroc/req/v3"
	"golang.org/x/net/html"
	hx "maragu.dev/gomponents-htmx/http"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
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
			return ui.ErrorPartial("Recipe not found"), nil
		}

		mainContent := ui.RecipeDetailPartial(recipe)

		recipes, err := s.GetRecipes(r.Context())

		// Second part updates another element out-of-band
		listContent := Div(
			ID("recipe-list"),
			Attr("hx-swap-oob", "true"), // Out-of-band swap
			ui.RecipeListPartial(recipes, id),
		)

		// Combine both parts in the response
		return Div(mainContent, listContent), nil

	}))

	// Get all recipes (main page)
	r.Get("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		recipes, err := s.GetRecipes(r.Context())
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
	r.Post("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		r.ParseForm()
		url := r.FormValue("url")

		resp := req.MustGet(url)

		doc, err := html.Parse(resp.Body)
		if err != nil {
			// ehh, maybe do nothing?
		}
		defer r.Body.Close()
		var title string
		var processAllProduct func(*html.Node) bool
		processAllProduct = func(n *html.Node) bool {
			if n.Type == html.ElementNode && n.Data == "title" {
				// process the Product details within each <li> element
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					// Extract the title text
					title = n.FirstChild.Data
					return true
				}

			}
			// traverse the child nodes
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				ret := processAllProduct(c)
				if ret {
					return true
				}
			}
			return false
		}
		// make a recursive call to your function
		processAllProduct(doc)

		id, err := s.AddRecipe(r.Context(), url, title, "")
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// Return just the updated list
		return ui.RecipeListPartial(recipes, id), nil
	}))

	r.Get("/recipes/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Show modal
		return ui.RecipeModal(), nil
	}))

	r.Post("/recipes/delete/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}

		err = s.DeleteRecipeByID(r.Context(), id)
		if err != nil {
			slog.Error("Could not delete recipe", "ID", id)
			return nil, err
		}
		slog.Info("Deleted recipe", "id", id)

		recipes, err := s.GetRecipes(r.Context())
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
	}))

	r.Get("/empty", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Clear modal
		return nil, nil
	}))
}
