package http

import (
	//"context"

	"net/http"
	"recipeze/ui"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"

	ghttp "maragu.dev/gomponents/http"
)

func RouteHome(r chi.Router) {
	r.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return ui.HomePage(ui.PageProps{}), nil
	}))

	r.Get("/login", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return ui.SignupForm("#modal-container"), nil
	}))
}
