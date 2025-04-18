package ui

import (
	. "maragu.dev/gomponents"
	"maragu.dev/gomponents-heroicons/v3/solid"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func CreateAccountPage(props PageProps) Node {
	props.Title = "Create account"

	return page(props,
		Div(
			Form(
				H2(Class("text-2xl font-bold text-gray-900 mb-6 text-center"),
					Text("Finish setting up account"),
				),
				//hx.Post(fmt.Sprintf("/recipes/update/%d", recipe.ID)),
				hx.Target("#recipe-detail"),

				// Editable Title
				Div(Class("mb-4"),
					Label(
						Class("block text-sm font-medium text-gray-700 mb-1"),
						For("recipe-name"),
						Text("Display name"),
					),
					Input(
						Type("text"),
						ID("display-name"),
						Name("name"),
						Placeholder("John Smith"),
						Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"),
					),
				),

				// Form Buttons
				Div(Class("flex justify-end space-x-3"),

					// Save Button
					Button(
						Type("submit"),
						Class("inline-flex justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"),
						Text("Save Changes"),
					),
				),
			),
		),
	)
}

func SignupForm(targetId string) Node {
	return Div(Class("max-w-md mx-auto bg-white p-8 rounded-xl shadow-md"),
		H2(Class("text-2xl font-bold text-gray-900 mb-6 text-center"),
			Text("Join Recipe Keeper"),
		),
		P(Class("text-gray-600 mb-6 text-center"),
			Text("Enter your email to get started. No passwords needed!"),
		),
		Form(
			ID("signup-form"),
			Class("space-y-4"),
			hx.Post("/auth/magic-link"),
			hx.Target(targetId),
			hx.Swap("innerHTML"),

			Div(Class("space-y-2"),
				Label(Class("block text-sm font-medium text-gray-700"), For("email"),
					Text("Email"),
				),
				Input(
					Type("email"),
					ID("email"),
					Name("email"),
					Class("w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"),
					Placeholder("you@email.com"),
					Required(),
				),
			),

			Button(
				Type("submit"),
				Class("w-full rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 cursor-pointer"),
				Text("Send Magic Link"),
			),

			P(Class("text-xs text-gray-500 mt-4 text-center"),
				Text("We'll email you a magic link for password-free sign in"),
			),
		),
	)
}

func SetupForm() Node {
	return Div(Class("max-w-md mx-auto bg-white p-8 rounded-xl shadow-md"),
		H2(Class("text-2xl font-bold text-gray-900 mb-6 text-center"),
			Text("Complete Your Account"),
		),
		P(Class("text-gray-600 mb-6 text-center"),
			Text("Set up your profile to start sharing recipes."),
		),
		Form(
			ID("setup-form"),
			Class("space-y-6"),
			hx.Post("/account/setup"),

			// Display Name
			Div(Class("space-y-2"),
				Label(Class("block text-sm font-medium text-gray-700"), For("display-name"),
					Text("Display Name"),
				),
				Input(
					Type("text"),
					ID("display-name"),
					Name("display_name"),
					Class("w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"),
					Placeholder("How you'll appear to others"),
					Required(),
				),
			),
			// Create First Group
			Div(Class("space-y-2 mt-8"),
				Div(Class("flex items-center"),
					solid.UserGroup(Class("h-5 w-5 text-indigo-500 mr-2")),
					H3(Class("text-lg font-medium text-gray-900"),
						Text("Create Your First Recipe Group"),
					),
				),
				P(Class("text-sm text-gray-600 ml-7"),
					Text("Recipe groups help you organize and share your recipes."),
				),

				// Group Name
				Div(Class("space-y-2 mt-4"),
					Label(Class("block text-sm font-medium text-gray-700"), For("group-name"),
						Text("Group Name"),
					),
					Input(
						Type("text"),
						ID("group-name"),
						Name("group_name"),
						Class("w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"),
						Placeholder("Family Favorites"),
						Value("My Recipes"),
					),
					P(Class("text-xs text-gray-500 mt-1"),
						Text("You can create more groups later"),
					),
				),
			),

			// Notification Preferences
			Div(Class("space-y-2 mt-8"),
				Div(Class("flex items-center"),
					solid.Bell(Class("h-5 w-5 text-indigo-500 mr-2")),
					H3(Class("text-lg font-medium text-gray-900"),
						Text("Notification Preferences"),
					),
				),

				// Email Notifications
				Div(Class("mt-4 space-y-3"),
					Div(Class("flex items-start"),
						Div(Class("flex items-center h-5"),
							Input(
								Type("checkbox"),
								ID("email-recipe-add"),
								Name("email_recipe_add"),
								Value("1"),
								Checked(),
								Class("h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"),
							),
						),
						Label(Class("ml-3 text-sm text-gray-700"), For("email-recipe-add"),
							Text("Email me when someone adds a recipe to a shared group"),
						),
					),
					Div(Class("flex items-start"),
						Div(Class("flex items-center h-5"),
							Input(
								Type("checkbox"),
								ID("email-group-invite"),
								Name("email_group_invite"),
								Value("1"),
								Checked(),
								Class("h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"),
							),
						),
						Label(Class("ml-3 text-sm text-gray-700"), For("email-group-invite"),
							Text("Email me when I'm invited to a recipe group"),
						),
					),
				),
			),

			// Invitation Section (Optional)
			Div(Class("space-y-2 mt-8 pt-6 border-t border-gray-200"),
				Div(Class("flex items-center"),
					solid.UserPlus(Class("h-5 w-5 text-indigo-500 mr-2")),
					H3(Class("text-lg font-medium text-gray-900"),
						Text("Invite Friends (Optional)"),
					),
				),
				P(Class("text-sm text-gray-600 ml-7"),
					Text("Share your recipes with family and friends."),
				),

				// Email Invitations
				Div(Class("space-y-2 mt-4"),
					Label(Class("block text-sm font-medium text-gray-700"), For("invite-emails"),
						Text("Email Addresses"),
					),
					Textarea(
						ID("invite-emails"),
						Name("invite_emails"),
						Class("w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"),
						Rows("3"),
						Placeholder("Enter email addresses, separated by commas"),
					),
					P(Class("text-xs text-gray-500 mt-1"),
						Text("We'll send them an invitation to join your recipe group"),
					),
				),

				// Personal Message
				Div(Class("space-y-2 mt-4"),
					Label(Class("block text-sm font-medium text-gray-700"), For("invite-message"),
						Text("Personal Message (Optional)"),
					),
					Textarea(
						ID("invite-message"),
						Name("invite_message"),
						Class("w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"),
						Rows("2"),
						Placeholder("Add a personal message to your invitation"),
					),
				),
			),

			// Submit Button
			Div(Class("pt-6"),
				Button(
					Type("submit"),
					Class("w-full rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 cursor-pointer"),
					Text("Complete Setup"),
				),
			),
		),
	)
}

// AccountSetupPage is the page wrapper for the setup form
func AccountSetupPage(props PageProps) Node {
	props.Title = "Complete Account Setup"
	return page(props,
		Div(Class("max-w-md mx-auto mt-8 mb-12"),
			SetupForm(),
		),
	)
}
