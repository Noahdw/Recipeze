package handler

// Handler is to make interacting with gomponents and requests more ergonomic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	rmiddleware "recipeze/middleware"
	"recipeze/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx/http"
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

func NewHandler(auth service.AuthService, recipe service.RecipeService) *handler {
	return &handler{
		AuthService:   auth,
		RecipeService: recipe,
	}
}

func InitRouting(r chi.Router, auth service.AuthService, recipe service.RecipeService) {
	mw := rmiddleware.NewAuthMiddleware(auth)
	h := NewHandler(auth, recipe)
	h.RouteHome(r, mw)
	h.RouteRecipe(r, mw)
}

func (h *handler) adapt(fn adaptFunc) http.HandlerFunc {
	return Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
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

func GetGroupID(r *http.Request) (int, error) {
	groupIDStr := chi.URLParam(r, "group_id")
	if groupIDStr == "" {
		return 0, fmt.Errorf("no group ID provided")
	}
	return strconv.Atoi(groupIDStr)
}

func getRecipeID(r *http.Request) (int, error) {
	groupIDStr := chi.URLParam(r, "recipe_id")
	if groupIDStr == "" {
		return 0, fmt.Errorf("no group ID provided")
	}
	return strconv.Atoi(groupIDStr)
}

type Handler = func(http.ResponseWriter, *http.Request) (Node, error)
type errorWithStatusCode interface {
	StatusCode() int
}

func Adapt(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n, err := h(w, r)
		if err != nil {

			if errors.Is(err, ErrDefault) {
				if hx.IsRequest(r.Header) {
					w.Header().Set("HX-Redirect", "/")
					w.WriteHeader(http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
				return
			}
			switch v := err.(type) {
			case errorWithStatusCode:
				w.WriteHeader(v.StatusCode())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		if n == nil {
			return
		}

		if err := n.Render(w); err != nil {
			http.Error(w, "error rendering node: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
