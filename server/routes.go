package server

import (
	"recipeze/handler"
	"recipeze/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"maragu.dev/httph"
)

func (s *server) SetupRoutes() {
	s.mux.Group(func(r chi.Router) {
		r.Use(middleware.Compress(5))

		// Sets up a static file handler with cache busting middleware.
		r.Group(func(r chi.Router) {
			r.Use(httph.VersionedAssets)
			setupStaticAssets(r)
		})

		recipeService := service.NewRecipeService(s.db)
		authService := service.NewAuthService(s.db)

		handler.InitRouting(r, authService, recipeService)
	})
}
