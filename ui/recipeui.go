package ui

import (
	"fmt"
	"log/slog"

	. "maragu.dev/gomponents"
	"maragu.dev/gomponents-heroicons/v3/solid"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"

	"recipeze/model"
)

// RecipePage shows the main recipe listing with a detail view
func RecipePage(props PageProps, recipes []model.Recipe, groupID int) Node {
	defaultId := 0
	var defaultRecipe *model.Recipe
	if len(recipes) > 0 {
		defaultId = recipes[0].ID
		defaultRecipe = &recipes[0]
	}
	slog.Info("id ", "id", defaultId)

	props.Title = "Recipes"

	return page(props,
		ModalContainer(),
		Div(Class("flex flex-col md:flex-row gap-6"),
			// Left column - Recipe List
			Div(Class("w-full md:w-1/3 "),
				H1(Class("text-2xl font-bold mb-4"), Text("Recipes")),
				AddRecipeButton(groupID),
				Div(ID("recipe-list"),
					RecipeListPartial(recipes, defaultId, groupID),
				),
			),
			// Right column - Recipe Detail
			Div(Class("w-full md:w-2/3 bg-gray-50 p-4 rounded-lg"),
				Div(ID("recipe-detail"),
					RecipeDetailPartial(defaultRecipe, groupID),
				),
			),
		),
	)
}

// RecipeListPartial shows the selectable recipe list
func RecipeListPartial(recipes []model.Recipe, selectedID int, groupID int) Node {
	if len(recipes) == 0 {
		return P(Text("No recipes yet. Add one to get started!"))
	}

	return Div(
		Class("max-h-[60vh] overflow-y-auto"), // Make scrollable
		Ul(Class("divide-y divide-gray-200"),
			Map(recipes, func(recipe model.Recipe) Node {
				return RecipeListItemPartial(&recipe, selectedID, groupID)
			}),
		),
	)
}

// RecipeListItemPartial shows an individual item in the larger recipe list
func RecipeListItemPartial(recipe *model.Recipe, selectedID int, groupID int) Node {
	var buttonClass string
	if recipe.ID == selectedID {
		buttonClass = "w-full text-left cursor-pointer bg-blue-100 hover:bg-blue-200 py-1 px-2 rounded font-medium active-recipe"
	} else {
		buttonClass = "w-full text-left cursor-pointer hover:bg-gray-100 py-1 px-2 rounded inactive-recipe"
	}
	id := fmt.Sprintf("recipe-list-item-%d", recipe.ID)
	return Li(
		Class("py-1"),
		Button(
			Class(buttonClass), ID(id),
			hx.Get(fmt.Sprintf("/g/%d/recipe/%d", groupID, recipe.ID)),
			hx.Target("#recipe-detail"),
			// Add class operations to clear previous selection
			Attr("hx-on::before-request", "document.querySelectorAll('.active-recipe').forEach(el => { el.classList.remove('bg-blue-100', 'hover:bg-blue-200', 'font-medium', 'active-recipe'); el.classList.add('hover:bg-gray-100', 'inactive-recipe'); })"),
			Text(recipe.Name),
		),
	)
}

// RecipeDetailPartial shows the details for a selected recipe
func RecipeDetailPartial(recipe *model.Recipe, groupID int) Node {
	if recipe == nil {
		return nil
	}
	return Div(
		H2(Class("text-xl font-bold mb-4"), Text(recipe.Name)), // title

		// Button container - flex row to make buttons appear horizontally
		Div(Class("flex flex-row gap-4 mb-6"),
			// View Original Recipe Link
			A(
				Attr("href", recipe.Url),
				Attr("target", "_blank"),
				Attr("rel", "noopener noreferrer"),
				Class("inline-flex items-center px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-md transition-colors"),
				Text("Go to recipe"),
			),

			// Edit Button
			Button(
				Class("inline-flex items-center justify-center rounded-md transition-colors bg-blue-500 hover:bg-blue-600 cursor-pointer"),
				hx.Get(fmt.Sprintf("/g/%d/recipes/update/%d", groupID, recipe.ID)),
				hx.Target("#recipe-detail"),
				Attr("aria-label", "Edit recipe"),
				Span(
					Class("flex items-center justify-center p-2"),
					solid.PencilSquare(Class("text-white h-5 w-5")),
				),
			),

			// Delete Button
			Button(
				Class("inline-flex items-center justify-center rounded-md transition-colors bg-red-300 hover:bg-red-600 cursor-pointer"),
				hx.Post(fmt.Sprintf("/g/%d/recipes/delete/%d", groupID, recipe.ID)),
				hx.Target("#recipe-detail"),
				Attr("aria-label", "Delete recipe"),
				Span(
					Class("flex items-center justify-center p-2"),
					solid.Trash(Class("text-white h-5 w-5")),
				),
			),
		),

		H3(Class("text-lg font-semibold mb-1"), Text("Notes")),
		Div(
			Class("whitespace-pre-wrap break-words"), // Preserves newlines and breaks long words
			Text(recipe.Description),
		),
		Img(
			Src(recipe.ImageURL),
			Class("w-full object-cover rounded-lg"),
			Loading("lazy"),
		),
	)
}

// RecipeEditPartial shows the details for a selected recipe in an editable form
func RecipeEditPartial(recipe *model.Recipe, groupID int) Node {
	return Div(
		Form(
			hx.Post(fmt.Sprintf("/g/%d/recipes/update/%d", groupID, recipe.ID)),
			hx.Target("#recipe-detail"),

			// Editable Title
			Div(Class("mb-4"),
				Label(
					Class("block text-sm font-medium text-gray-700 mb-1"),
					For("recipe-name"),
					Text("Recipe Name"),
				),
				Input(
					Type("text"),
					ID("recipe-name"),
					Name("name"),
					Value(recipe.Name),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"),
				),
			),

			// Editable URL
			Div(Class("mb-4"),
				Label(
					Class("block text-sm font-medium text-gray-700 mb-1"),
					For("recipe-url"),
					Text("Recipe URL"),
				),
				Input(
					Type("url"),
					ID("recipe-url"),
					Name("url"),
					Value(recipe.Url),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"),
				),
			),

			// Editable Notes
			Div(Class("mb-4"),
				Label(
					Class("block text-sm font-medium text-gray-700 mb-1"),
					For("recipe-description"),
					Text("Notes"),
				),
				Textarea(
					ID("recipe-description"),
					Name("description"),
					Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"),
					Text(recipe.Description),
					Rows("6"),
				),
			),

			// Form Buttons
			Div(Class("flex justify-end space-x-3"),
				// Cancel Button
				Button(
					Type("button"),
					Class("px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"),
					hx.Get(fmt.Sprintf("/g/%d/recipe/%d", groupID, recipe.ID)),
					hx.Target("#recipe-detail"),
					Text("Cancel"),
				),

				// Save Button
				Button(
					Type("submit"),
					Class("inline-flex justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"),
					Text("Save Changes"),
				),
			),
		),
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

func AddRecipeButton(group_id int) Node {
	return Button(
		Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer"),
		hx.Get(fmt.Sprintf("/g/%d/recipes/new", group_id)),
		hx.Target("#modal-container"),
		hx.Swap("innerHTML"),
		// focus the URL input after the modal is loaded
		Attr("hx-on::after-request", "setTimeout(() => document.getElementById('recipe-url').focus(), 10)"),
		Text("Add Recipe"),
	)
}

func ModalContainer() Node {
	return Div(
		ID("modal-container"),
		// The modal will be loaded here
	)
}

func RecipeModal(group_id int) Node {
	return Div(
		Class("fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50"),
		Div(
			Class("bg-white rounded-lg p-6 max-w-md w-full"),
			Div(
				Class("flex justify-between items-center mb-4"),
				H3(Class("text-lg font-medium"), Text("Add New Recipe")),
				Button(
					Class("text-gray-400 hover:text-gray-500 cursor-pointer"),
					hx.Get("/empty"),
					hx.Target("#modal-container"),
					hx.Swap("innerHTML"),
					Text("×"),
				),
			),
			Form(
				hx.Post(fmt.Sprintf("/g/%d/recipes", group_id)),
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
						Class("mr-3 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded cursor-pointer"),
						hx.Get("/empty"),
						hx.Target("#modal-container"),
						hx.Swap("innerHTML"),
						Text("Cancel"),
					),
					Button(
						Type("submit"),
						Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer"),
						Text("Save Recipe"),
					),
				),
			),
		),
	)
}
