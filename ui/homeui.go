package ui

//"time"
import (
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

// HomePage is the front page of the app
func HomePage(props PageProps) Node {
	props.Title = "Home"
	targetId := "#modal-container"

	return page(props,
		ModalContainer(),
		// Hero section
		HeroSection(targetId),

		// How it works
		HowItWorksSection(),

		// Popular categories
		CategoriesSection(),

		// Testimonial
		TestimonialSection(),
	)
}

func PrimaryButton(text string, hxGet string, hxTarget string) Node {
	return Button(
		Class("rounded-xl bg-indigo-600 px-6 py-3 text-lg font-medium text-white shadow-md hover:bg-indigo-500 transition-colors duration-200 cursor-pointer"),
		Text(text),
		hx.Get(hxGet),
		hx.Target(hxTarget),
	)
}

func SecondaryButton(text string, hxGet string, hxTarget string) Node {
	return Button(
		Class("rounded-xl bg-white px-6 py-3 text-lg font-medium text-indigo-600 border-2 border-indigo-200 shadow-md hover:border-indigo-300 hover:bg-indigo-50 transition-colors duration-200"),
		Text(text),
		hx.Get(hxGet),
		hx.Target(hxTarget),
	)
}

// Section container with consistent styling
func SectionContainer(bgColor string, className string, children ...Node) Node {
	// Create attributes for the outer div
	outerAttrs := []Node{Class("py-16 " + bgColor + " " + className)}

	// Create attributes for the inner div
	innerAttrs := []Node{Class("max-w-7xl mx-auto px-6")}

	// Add children to inner attributes
	for _, child := range children {
		innerAttrs = append(innerAttrs, child)
	}

	// Create the inner div and add it to outer attributes
	outerAttrs = append(outerAttrs, Div(innerAttrs...))

	// Return the composed divs
	return Div(outerAttrs...)
}

// Rounded section with consistent styling
func RoundedSection(bgColor string, children ...Node) Node {
	return SectionContainer(bgColor, "rounded-3xl mx-4 my-8 md:mx-8", children...)
}

// Section header with consistent styling
func SectionHeader(text string) Node {
	return H2(Class("text-3xl font-bold text-center text-gray-900"), Text(text))
}

// Step item component for "How it works" section
func StepItem(number string, title string, description string) Node {
	return Div(Class("flex flex-col items-center text-center"),
		Div(Class("rounded-2xl bg-indigo-100 p-4 h-16 w-16 flex items-center justify-center text-xl font-bold text-indigo-600"),
			Text(number),
		),
		H3(Class("mt-5 text-xl font-semibold text-gray-900"),
			Text(title),
		),
		P(Class("mt-3 text-gray-600"),
			Text(description),
		),
	)
}

// Category card component
func CategoryCard(text string, href string) Node {
	return A(
		Class("bg-indigo-100 rounded-2xl p-8 flex items-center justify-center text-lg font-medium hover:bg-indigo-200 transition-colors duration-200 shadow-sm"),
		Href(href),
		Text(text),
	)
}

// Hero section component
func HeroSection(targetId string) Node {
	return Div(Class("flex flex-col items-center text-center py-16 px-6 max-w-4xl mx-auto"), ID(targetId),
		H1(Class("text-4xl font-bold text-gray-900"),
			Text("Save your favorite recipes in one click"),
		),
		P(Class("mt-6 text-xl text-gray-600"),
			Text("No more lost bookmarks or screenshots. Just grab recipes you love and share them with your family."),
		),
		Div(Class("mt-10 flex flex-wrap gap-4 justify-center"),
			PrimaryButton("Get cooking", "/login", targetId),
		),
	)
}

// How it works section component
func HowItWorksSection() Node {
	return RoundedSection("bg-gray-50",
		SectionHeader("How it works (it's super simple)"),
		Div(Class("mt-12 grid grid-cols-1 gap-10 md:grid-cols-3"),
			StepItem("1", "Find a tasty recipe", "See something yummy online? Just copy the URL and we'll save it."),
			StepItem("2", "Invite your family", "Got family recipes to share? Invite everyone to contribute their favorites."),
			StepItem("3", "Cook & enjoy", "Access your recipes anywhere - perfect for grocery shopping or kitchen disasters."),
		),
	)
}

// Categories section component
func CategoriesSection() Node {
	return SectionContainer("", "",
		SectionHeader("What's cooking?"),
		Div(Class("mt-12 grid grid-cols-2 gap-6 md:grid-cols-4"),
			CategoryCard("Breakfast Faves", "/recipes?category=breakfast"),
			CategoryCard("Dinner Ideas", "/recipes?category=dinner"),
			CategoryCard("Sweet Treats", "/recipes?category=desserts"),
			CategoryCard("Holiday Magic", "/recipes?category=holidays"),
		),
	)
}

// Testimonial section component
func TestimonialSection() Node {
	return Div(Class("py-16 bg-indigo-50 rounded-3xl mx-4 mb-12 md:mx-8"),
		Div(Class("max-w-3xl mx-auto px-6 text-center"),
			SectionHeader("People are loving it"),
			P(Class("mt-10 italic text-xl text-gray-700"),
				Text("\"Finally found a way to save Grandma's secret sauce recipe where the whole family can access it. No more texting 'hey, how much garlic goes in that thing again?'\""),
			),
			P(Class("mt-4 font-medium text-gray-900"),
				Text("— The Johnson Family"),
			),
		),
	)
}
