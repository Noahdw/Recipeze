package handler

// Handler is to make interacting with gomponents and requests more ergonomic

import (
	"context"
	"net/http"
	"recipeze/service"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

type adaptFunc = func(ctx requestContext) (Node, error)

type requestContext struct {
	w http.ResponseWriter
	r *http.Request
}

type handler struct {
	service.AuthService
	service.RecipeService
}

func SetupRouting(r chi.Router, auth service.AuthService, recipe service.RecipeService) {
	h := &handler{
		AuthService:   auth,
		RecipeService: recipe,
	}
	h.RouteHome(r)
	h.RouteRecipe(r)
}

func (h *handler) adapt(fn adaptFunc) http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		ctx := requestContext{
			w: w,
			r: r,
		}
		return fn(ctx)
	})
}

func (c *requestContext) context() context.Context {
	return c.r.Context()
}

func (c *requestContext) queryParam(param string) string {
	return c.r.URL.Query().Get(param)
}
