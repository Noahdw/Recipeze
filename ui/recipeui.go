package ui

import (
	"fmt"
	"log/slog"

	. "maragu.dev/gomponents"
	"maragu.dev/gomponents-heroicons/v3/solid"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"

	"recipeze/model"
	"recipeze/parsing"
)

// RecipePage shows the main recipe listing with a detail view
func RecipePage(props PageProps, recipes []model.Recipe, group model.Group) Node {
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

		// Group Indicator and Member Management
		Div(Class("flex items-center justify-between gap-2 mb-2"),
			Div(Class(""),
				AddRecipeButton(group.ID),
				AddPlanMealsButton(group.ID),
			),
			Div(Class("flex items-center gap-2"),
				// Group Selector Dropdown
				Div(Class("relative inline-block"),
					Button(
						ID("group-selector"),
						Class("flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white px-3 py-2 rounded"),
						Attr("aria-haspopup", "true"),
						Attr("aria-expanded", "false"),
						Attr("onclick", "document.getElementById('group-dropdown').classList.toggle('hidden')"),
						solid.UserGroup(Class("h-5 w-5")),
						Text(group.Name),
						solid.ChevronDown(Class("h-4 w-4")),
					),
					groupDropDownMenu(&group),
				),
				groupMembersDisplay(&group),
			),
		),

		Div(Class("flex flex-col md:flex-row gap-6"),
			// Left column - Recipe List
			Div(Class("w-full md:w-1/3"),
				H1(Class("text-2xl font-bold mb-4"), Text("Recipes")),
				Div(ID("recipe-list"),
					RecipeListPartial(recipes, defaultId, group.ID),
				),
			),
			// Right column - Recipe Detail
			Div(Class("w-full md:w-2/3 bg-gray-50 p-4 rounded-lg"),
				Div(ID("recipe-detail"),
					RecipeDetailPartial(defaultRecipe, group.ID),
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
		buttonClass = "w-full text-left cursor-pointer bg-blue-100 hover:bg-blue-200 py-1 px-2 rounded active-recipe"
	} else {
		buttonClass = "w-full text-left cursor-pointer hover:bg-gray-100 py-1 px-2 rounded inactive-recipe"
	}
	id := fmt.Sprintf("recipe-list-item-%d", recipe.ID)
	return Li(
		Button(
			Class(buttonClass), ID(id),
			hx.Get(fmt.Sprintf("/g/%d/recipe/%d", groupID, recipe.ID)),
			hx.Target("#recipe-detail"),
			// Add class operations to clear previous selection
			Attr("hx-on::before-request", "document.querySelectorAll('.active-recipe').forEach(el => { el.classList.remove('bg-blue-100', 'hover:bg-blue-200', 'active-recipe'); el.classList.add('hover:bg-gray-100', 'inactive-recipe'); })"),
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
		H3(Class("text-lg font-semibold mb-1"), Text("Ingredients")),
		Div(
			Class("whitespace-pre-wrap break-words"), // Preserves newlines and breaks long words
			Text(parsing.RecipeIngredients(recipe.Data)),
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
		Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer mr-2"),
		hx.Get(fmt.Sprintf("/g/%d/recipes/new", group_id)),
		hx.Target("#modal-container"),
		hx.Swap("innerHTML"),
		// focus the URL input after the modal is loaded
		Attr("hx-on::after-request", "setTimeout(() => document.getElementById('recipe-url').focus(), 10)"),
		Text("Add Recipe"),
	)
}

func AddInviteButton(group_id int) Node {
	return Button(
		Class("hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer mr-2"),
		hx.Get(fmt.Sprintf("/g/%d/recipes/invite", group_id)),
		hx.Target("#modal-container"),
		hx.Swap("innerHTML"),
		// focus the URL input after the modal is loaded
		Attr("hx-on::after-request", "setTimeout(() => document.getElementById('recipe-invite').focus(), 10)"),
		Text("Invite"),
	)
}

func AddPlanMealsButton(group_id int) Node {
	return Button(
		Class("bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded cursor-pointer"),
		hx.Get(fmt.Sprintf("/g/%d/recipes/new", group_id)),
		hx.Target("#modal-container"),
		hx.Swap("innerHTML"),
		// focus the URL input after the modal is loaded
		Attr("hx-on::after-request", "setTimeout(() => document.getElementById('recipe-url').focus(), 10)"),
		Text("Plan meals"),
	)
}

func ModalContainer() Node {
	return Div(
		ID("modal-container"),
		// The modal will be loaded here
	)
}

func groupDropDownMenu(group *model.Group) Node {
	return Div(
		ID("group-dropdown"),
		Class("hidden absolute left-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-10"),
		Div(Class("py-1"),
			Div(
				Class("px-4 py-2 text-xs text-gray-500"),
				Text("YOUR GROUPS"),
			),
			// loop through groups here
			A(
				Href(fmt.Sprintf("/g/%d/recipes", group.ID)),
				Class("block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"),
				Div(Class("flex items-center justify-between"),
					Div(Class("flex items-center gap-2"),
						solid.UserGroup(Class("h-4 w-4 text-gray-400")),
						Text(group.Name),
					),
					Span(Class("text-xs text-gray-500"),
						Text(fmt.Sprintf("%d members", len(group.Members))),
					),
				),
			),
			// Create New Group option
			Div(Class("border-t border-gray-100 mt-1 pt-1"),
				Button(
					Class("w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"),
					hx.Get("/groups/new"),
					hx.Target("#modal-container"),
					Div(Class("flex items-center gap-2"),
						solid.Plus(Class("h-4 w-4 text-green-500")),
						Text("Create New Group"),
					),
				),
			),
		),
	)
}

func groupMembersDisplay(group *model.Group) Node {
	return Div(
		Class("ml-4 flex items-center"),
		// Member avatars
		Div(Class("flex -space-x-2"),
			// We'd iterate through the first few members here
			Map(group.Members[:min(3, len(group.Members))], func(member model.GroupMember) Node {
				return Div(
					Class("w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-xs border-2 border-white"),
					Text(member.Name[:1]), // First letter of name as avatar
				)
			}),
			// Show +X more if there are more members
			If(len(group.Members) > 3,
				Div(
					Class("w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-xs border-2 border-white"),
					Text(fmt.Sprintf("+%d", len(group.Members)-3)),
				),
			),
		),
		// Add member button
		Button(
			Class("ml-2 w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center hover:bg-blue-200"),
			hx.Get(fmt.Sprintf("/g/%d/members/invite", group.ID)),
			hx.Target("#modal-container"),
			Attr("aria-label", "Invite group member"),
			solid.Plus(Class("h-4 w-4")),
		),
		// View all members button
		Button(
			Class("ml-2 text-sm text-blue-600 hover:text-blue-800"),
			hx.Get(fmt.Sprintf("/g/%d/members", group.ID)),
			hx.Target("#modal-container"),
			Text("View all"),
		),
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

func InviteModal(group_id int) Node {
	return Div(
		Class("fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50"),
		Div(
			Class("bg-white rounded-lg p-6 max-w-md w-full"),
			Div(
				Class("flex justify-between items-center mb-4"),
				H3(Class("text-lg font-medium"), Text("Invite Someone")),
				Button(
					Class("text-gray-400 hover:text-gray-500 cursor-pointer"),
					hx.Get("/empty"),
					hx.Target("#modal-container"),
					hx.Swap("innerHTML"),
					Text("×"),
				),
			),
			Form(
				hx.Post(fmt.Sprintf("/g/%d/recipes/invite", group_id)),
				hx.Target("#recipe-list"),
				hx.Swap("innerHTML"),
				Attr("hx-on::after-request", "document.querySelector('#modal-container').innerHTML = ''"),

				Div(
					Class("mb-4"),
					Label(Class("block text-sm font-medium text-gray-700"), For("recipe-invite"), Text("Enter the email to receive the invite")),
					Input(Type("email"), ID("recipe-invite"), Name("email"), Required(), Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm")),
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
						Text("Send invite"),
					),
				),
			),
		),
	)
}
