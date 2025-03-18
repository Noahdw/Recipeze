package handler

import (
	"recipeze/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"maragu.dev/httph"
)

func (s *Server) setupRoutes() {
	s.mux.Group(func(r chi.Router) {
		r.Use(middleware.Compress(5))

		// Sets up a static file handler with cache busting middleware.
		r.Group(func(r chi.Router) {
			r.Use(httph.VersionedAssets)

			Static(r)
		})

		recipeService := service.NewRecipeService(s.db)
		authService := service.NewAuthService(s.db)
		handler := NewHandler(authService, recipeService)
		RouteRecipe(r, handler)
		RouteHome(r, handler)
	})
}
