package ui

import (
	. "maragu.dev/gomponents"
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
					Placeholder("your@email.com"),
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
