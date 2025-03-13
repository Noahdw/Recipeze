package http

import (
	//"context"

	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	"recipeze/model"
	"recipeze/repo"
	"recipeze/service"
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

func RouteRecipe(r chi.Router, s *service.Recipe) {
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

		meta := extractMeta(doc)
		id, err := s.AddRecipe(r.Context(), url, meta.title, meta.description, meta.imageURL)
		if err != nil {
			slog.Error("Could not add recipe", "error", err.Error())
			return nil, err
		}

		recipes, err := s.GetRecipes(r.Context())
		if err != nil {
			slog.Error("Could not get recipes", "error", err.Error())
			return nil, err
		}

		recipe, err := s.GetRecipeByID(r.Context(), int32(id))

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

	r.Get("/recipes/update/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		idstr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, err
		}
		recipe, err := s.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return nil, err
		}

		return ui.RecipeEditPartial(recipe), nil
	}))

	r.Post("/recipes/update/{id}", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
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
		err = s.UpdateRecipe(r.Context(), repo.UpdateRecipeParams{
			ID:          int32(id),
			Name:        repo.StringPG(r.FormValue("name")),
			Url:         repo.StringPG(r.FormValue("url")),
			Description: repo.StringPG(r.FormValue("description")),
		})
		if err != nil {
			slog.Error("Could not update recipe", "ID", id)
			return nil, err
		}

		recipe, err := s.GetRecipeByID(r.Context(), int32(id))
		if err != nil {
			return nil, err
		}

		mainContent := ui.RecipeDetailPartial(recipe)
		recipes, err := s.GetRecipes(r.Context())
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
	}))

	r.Get("/empty", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Clear modal
		return nil, nil
	}))

}

// extractMetaImage finds the og:image or similar meta tag from an HTML document
func extractMeta(doc *html.Node) meta {
	var m meta
	// Try to find Open Graph image first
	ogImage := findMetaContent(doc, "property", "og:image")
	if ogImage != "" {
		m.imageURL = ogImage
	}

	// Try Twitter image
	twitterImage := findMetaContent(doc, "name", "twitter:image")
	if twitterImage != "" {
		m.imageURL = twitterImage
	}

	// Try regular meta image
	metaImage := findMetaContent(doc, "name", "image")
	if metaImage != "" {
		m.imageURL = metaImage
	}

	m.title = findMetaContent(doc, "property", "og:title")
	if m.title == "" {
		// TODO: Look for non meta title
		m.title = "recipe"
	}
	m.description = findMetaContent(doc, "name", "description")
	m.siteName = findMetaContent(doc, "property", "og:site_name")
	return m
}

// findMetaContent looks for a meta tag with the specified attribute and value
// and returns its content attribute
func findMetaContent(doc *html.Node, attrName, attrValue string) string {
	var content string

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			// Check if this meta tag has the attribute we're looking for
			var hasAttr bool
			var contentAttr string

			for _, attr := range n.Attr {
				if attr.Key == attrName && attr.Val == attrValue {
					hasAttr = true
				}
				if attr.Key == "content" {
					contentAttr = attr.Val
				}
			}

			if hasAttr && contentAttr != "" {
				content = contentAttr
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if content == "" { // Only continue if we haven't found it yet
				crawler(c)
			}
		}
	}

	crawler(doc)
	return content
}
