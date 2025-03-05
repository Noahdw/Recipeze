package http

import (
	//"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	"recipeze/service"
	"recipeze/ui"

	"github.com/imroc/req/v3"
	"golang.org/x/net/html"
	hx "maragu.dev/gomponents-htmx/http"
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

		return ui.RecipeDetailPartial(recipe), nil
	}))

	// Get all recipes (main page)
	r.Get("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// If HTMX request, return just the list
		if hx.IsRequest(r.Header) {
			return ui.RecipeListPartial(recipes, time.Now()), nil
		}

		// Otherwise return full page
		return ui.RecipePage(ui.PageProps{}, recipes, time.Now()), nil
	}))

	// Add new recipe
	r.Post("/recipes", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		r.ParseForm()
		url := r.FormValue("url")
		req.DevMode()

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

		err = s.AddRecipe(r.Context(), url, title, "")
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			return nil, err
		}

		// Return just the updated list
		return ui.RecipeListPartial(recipes, time.Now()), nil
	}))

	r.Get("/recipes/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Show modal
		return ui.RecipeModal(), nil
	}))

	r.Get("/empty", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Clear modal
		return nil, nil
	}))
}
