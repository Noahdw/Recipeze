package html

import (
	"fmt"
	"time"

	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"

	"recipeze/model"
)

// RecipePage shows the main recipe listing with a detail view
func RecipePage(props PageProps, recipes []model.Recipe, now time.Time) Node {
	props.Title = "Recipes"
	return page(props,
		Div(Class("flex flex-col md:flex-row gap-6"),
			// Left column - Recipe List
			Div(Class("w-full md:w-1/3"),
				H1(Class("text-2xl font-bold mb-4"), Text("Recipes")),
				Button(
					Class("rounded-md bg-indigo-600 px-2.5 py-1.5 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 mb-4"),
					Text("Add recipe"), hx.Post("/recipes"), hx.Target("#recipe-list")),
				Div(ID("recipe-list"),
					RecipeListPartial(recipes, now),
				),
			),
			// Right column - Recipe Detail
			Div(Class("w-full md:w-2/3 bg-gray-50 p-4 rounded-lg"),
				Div(ID("recipe-detail"),
					// Initially empty or with a placeholder
					P(Class("text-gray-500 italic"), Text("Select a recipe to view details")),
				),
			),
		),
	)
}

// RecipeListPartial shows the selectable recipe list
func RecipeListPartial(recipes []model.Recipe, now time.Time) Node {
	if len(recipes) == 0 {
		return P(Text("No recipes yet. Add one to get started!"))
	}

	return Ul(Class("divide-y divide-gray-200"),
		Map(recipes, func(recipe model.Recipe) Node {
			return Li(
				Class("py-2 px-3"),
				Button(
					Class("w-full text-left cursor-pointer hover:bg-gray-100 py-1 px-2 rounded"),
					hx.Get(fmt.Sprintf("/recipes/%d", recipe.ID)),
					hx.Target("#recipe-detail"),
					Text(recipe.Name),
				),
			)
		}),
	)
}

// RecipeDetailPartial shows the details for a selected recipe
func RecipeDetailPartial(recipe *model.Recipe) Node {
	return Div(
		H2(Class("text-xl font-bold mb-4"), Text(recipe.Name)),
		P(Class("mb-2"),
			A(Href(recipe.Url), Target("_blank"), Text("View Original Recipe")),
		),
		H3(Class("text-lg font-semibold mb-2"), Text("Notes")),
		P(Text(recipe.Url)),
		// Add more recipe details here as needed
	)
}

// ErrorPartial displays an error message to the user
func ErrorPartial(message string) Node {
	return Div(
		Class("bg-red-50 border border-red-200 text-red-800 rounded-md p-4"),
		Div(
			Class("flex"),
			Div(
				Class("flex-shrink-0"),
				// You could add an SVG icon here if desired
			),
			Div(
				Class("ml-3"),
				P(
					Class("text-sm font-medium text-red-800"),
					Text(message),
				),
			),
		),
	)
}
