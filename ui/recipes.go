package ui

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
		ModalContainer(),
		Div(Class("flex flex-col md:flex-row gap-6"),
			// Left column - Recipe List
			Div(Class("w-full md:w-1/3"),
				H1(Class("text-2xl font-bold mb-4"), Text("Recipes")),
				AddRecipeButton(),
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

func AddRecipeButton() Node {
	return Button(
		Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"),
		hx.Get("/recipes/new"),
		hx.Target("#modal-container"),
		hx.Swap("innerHTML"),
		Text("Add Recipe"),
	)
}

func ModalContainer() Node {
	return Div(
		ID("modal-container"),
		// The modal will be loaded here
	)
}

func RecipeModal() Node {
	return Div(
		Class("fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50"),
		Div(
			Class("bg-white rounded-lg p-6 max-w-md w-full"),
			Div(
				Class("flex justify-between items-center mb-4"),
				H3(Class("text-lg font-medium"), Text("Add New Recipe")),
				Button(
					Class("text-gray-400 hover:text-gray-500"),
					hx.Get("/empty"),
					hx.Target("#modal-container"),
					hx.Swap("innerHTML"),
					Text("×"),
				),
			),
			Form(
				hx.Post("/recipes"),
				hx.Target("#recipe-list"),
				hx.Swap("innerHTML"),
				//hx.On("after-request", "document.querySelector('#modal-container').innerHTML = ''"),
				Attr("hx-on::after-request", "document.querySelector('#modal-container').innerHTML = ''"),

				Div(
					Class("mb-4"),
					Label(Class("block text-sm font-medium text-gray-700"), For("recipe-url"), Text("Recipe URL")),
					Input(Type("url"), ID("recipe-url"), Name("url"), Required(), Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm")),
				),

				Div(
					Class("mt-6 flex justify-end"),
					Button(
						Type("button"),
						Class("mr-3 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded"),
						hx.Get("/empty"),
						hx.Target("#modal-container"),
						hx.Swap("innerHTML"),
						Text("Cancel"),
					),
					Button(
						Type("submit"),
						Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"),
						Text("Save Recipe"),
					),
				),
			),
		),
	)
}
