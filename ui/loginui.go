package ui

import (
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

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
				Class("w-full rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"),
				Text("Send Magic Link"),
			),

			P(Class("text-xs text-gray-500 mt-4 text-center"),
				Text("We'll email you a magic link for password-free sign in"),
			),
		),
	)
}
