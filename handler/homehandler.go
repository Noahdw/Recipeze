package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"recipeze/appconfig"

	mw "recipeze/middleware"
	"recipeze/ui"

	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	. "maragu.dev/gomponents"
)

type ErrType int

var (
	ErrDefault       = errors.New("unknown error")
	ErrNotFound      = errors.New("URL not found")
	ErrNotAuthorized = errors.New("not auhotirzed")
)

func (h *handler) RouteHome(r chi.Router, m *mw.AuthMiddleware) {

	r.Get("/", h.adapt(func(ctx requestContext) (Node, error) {
		return ui.HomePage(ui.PageProps{IncludeHeader: false}), nil
	}))
	r.Get("/login", h.login())
	r.Get("/logout", h.logout())
	r.Post("/auth/magic-link", h.sendMagicLinkToEmail())
	r.Get("/auth/verify", h.authenticateByMagicLink())

	r.Group(func(r chi.Router) {
		r.Use(m.Authenticate)
		r.Get("/account/setup", h.showSetupUserDetails())
		r.Post("/account/setup", h.finishSetupUserDetails())
	})
}

func (h *handler) login() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(ctx.r, "session")
		session_token, ok := session.Values["session_token"]
		if !ok {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, nil
		}

		sessionTokenStr, ok := session_token.(string)
		if !ok {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, nil
		}

		user, err := h.GetLoggedInUser(ctx.context(), sessionTokenStr)
		if err != nil {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, nil
		}

		groups, err := h.GetUserGroups(ctx.context(), user.ID)
		if err != nil {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, nil
		}

		// Redirect to the default recipes page for the user
		url := fmt.Sprintf("%s/g/%d/recipes", appconfig.Config.URL, groups[0].ID)
		ctx.w.Header().Set("HX-Redirect", url)
		ctx.w.WriteHeader(http.StatusOK)
		return nil, nil
	})
}

func (h *handler) logout() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(ctx.r, "session")
		session.Values["session_token"] = ""
		err := session.Save(ctx.r, ctx.w)
		if err != nil {
			return nil, ErrDefault
		}
		ctx.w.Header().Set("HX-Redirect", appconfig.Config.URL)
		ctx.w.WriteHeader(http.StatusOK)
		return nil, nil
	})
}

func (h *handler) sendMagicLinkToEmail() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {

		config, err := awsc.LoadDefaultConfig(ctx.context(),
			awsc.WithRegion("us-east-2"),
		)

		if err != nil {
			// TODO
			slog.Error("Could not load email config")
			return nil, ErrDefault
		}

		// Get email user entered
		err = ctx.r.ParseForm()
		if err != nil {
			return nil, ErrDefault
		}
		email := ctx.r.FormValue("email")

		destination := &types.Destination{
			ToAddresses: []string{email},
		}

		// Use users email to create a auth token
		// Auth token will allow use to verify their email
		token, err := h.CreateRegistrationToken(ctx.context(), email)
		if err != nil {
			return nil, ErrDefault
		}

		// Construct the magic link to be used for verification
		magicLink := appconfig.Config.URL + "/auth/verify?token=" + token

		params := &sesv2.SendEmailInput{
			Content:              createLoginEmail(appconfig.AppName(), magicLink),
			ConfigurationSetName: new(string),
			Destination:          destination,
			FromEmailAddress:     &appconfig.Config.FromEmail,
		}

		client := sesv2.NewFromConfig(config)
		_, err = client.SendEmail(ctx.context(), params)
		if err != nil {
			slog.Error("Could not send email", "to", email, "err", err.Error())
			return nil, ErrDefault
		}

		slog.Info("sending magic link", "email", email)
		return ui.SignupForm("#modal-container"), nil
	})
}

func (h *handler) authenticateByMagicLink() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		token := ctx.queryParam("token")
		email, err := h.VerifyRegistrationToken(ctx.context(), token, ctx.r)
		if err != nil {
			slog.Error("Could not verify email", "err", err.Error())
			url := appconfig.Config.URL
			http.Redirect(ctx.w, ctx.r, url, http.StatusSeeOther)
			return nil, ErrDefault
		}

		slog.Info("Verified email")

		user, err := h.GetUser(ctx.context(), email)
		if err != nil {
			// Assume account does not exist
			//return ui.CreateAccountPage(ui.PageProps{}), nil
			user, err = h.CreateAccount(ctx.context(), email)
			if err != nil {
				return nil, ErrDefault
			}

		}

		session_token := generateSessionToken(32)
		loggedIn, err := h.Login(ctx.context(), user.ID, session_token)
		if err != nil {
			return nil, ErrDefault
		}
		if !loggedIn {
			return nil, ErrDefault
		}

		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(ctx.r, "session")
		// Set some session values.
		session.Values["session_token"] = session_token
		err = session.Save(ctx.r, ctx.w)
		if err != nil {
			return nil, ErrDefault
		}

		url := appconfig.Config.URL + "/account/setup"
		http.Redirect(ctx.w, ctx.r, url, http.StatusSeeOther)
		slog.Info("Redirect verified user", "url", url)

		return nil, nil
	})
}

func (h *handler) showSetupUserDetails() http.HandlerFunc {
	slog.Info("requested showSetupUserDetails")
	return h.adapt(func(ctx requestContext) (Node, error) {
		user := mw.GetUserFromContext(ctx.context())
		if user == nil {
			return nil, ErrDefault
		}
		if user.SetupComplete {

			return nil, ErrDefault
		}
		props := ui.PageProps{
			Title:         "",
			Description:   "",
			IncludeHeader: false,
			GroupID:       0,
		}
		return ui.AccountSetupPage(props), nil
	})
}

func (h *handler) finishSetupUserDetails() http.HandlerFunc {
	slog.Info("requested finishSetupUserDetails")
	return h.adapt(func(ctx requestContext) (Node, error) {
		user := mw.GetUserFromContext(ctx.context())
		if user == nil {
			return nil, ErrDefault
		}
		groups, err := h.GetUserGroups(ctx.context(), user.ID)
		if err != nil {
			return nil, ErrDefault
		}

		// Redirect to the default recipes page for the user
		url := fmt.Sprintf("%s/g/%d/recipes", appconfig.Config.URL, groups[0].ID)
		ctx.w.Header().Set("HX-Redirect", url)
		ctx.w.WriteHeader(http.StatusOK)
		return nil, nil
	})
}

func generateSessionToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func renderNode(w http.ResponseWriter, r *http.Request, node Node) {
	Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return node, nil
	})(w, r)
}
