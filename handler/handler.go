package handler

// Handler is to make interacting with gomponents and requests more ergonomic

import (
	"context"
	"net/http"
	"recipeze/service"

	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

type AdaptFunc = func(ctx RequestContext) (Node, error)

type RequestContext struct {
	w http.ResponseWriter
	r *http.Request
}

type handler struct {
	service.AuthService
	service.RecipeService
}

func NewHandler(auth service.AuthService, recipe service.RecipeService) *handler {
	return &handler{
		AuthService:   auth,
		RecipeService: recipe,
	}
}

func (h *handler) Adapt(fn AdaptFunc) http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		ctx := RequestContext{
			w: w,
			r: r,
		}
		return fn(ctx)
	})
}

func (c *RequestContext) Context() context.Context {
	return c.r.Context()
}

func (c *RequestContext) QueryParam(param string) string {
	return c.r.URL.Query().Get(param)
}
