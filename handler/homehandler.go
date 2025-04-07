package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"recipeze/appconfig"

	"recipeze/ui"

	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func (h *handler) RouteHome(r chi.Router) {

	r.Get("/", h.adapt(func(ctx requestContext) (Node, error) {
		return ui.HomePage(ui.PageProps{IncludeHeader: false}), nil
	}))
	r.Get("/login", h.login())

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(r, "session")
		session.Values["session_token"] = ""
		err := session.Save(r, w)
		if err != nil {
			return
		}
		w.Header().Set("HX-Redirect", appconfig.Config.URL)
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/auth/magic-link", h.sendMagicLinkToEmail())

	r.Get("/auth/verify", h.authenticateByMagicLink())
}

func (h *handler) login() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(ctx.r, "session")
		session_token, ok := session.Values["session_token"]
		if !ok {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, fmt.Errorf("")
		}

		sessionTokenStr, ok := session_token.(string)
		if !ok {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, fmt.Errorf("")
		}

		user, err := h.GetLoggedInUser(ctx.context(), sessionTokenStr)
		if err != nil {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, err
		}

		groups, err := h.GetUserGroups(ctx.context(), user.ID)
		if err != nil {
			renderNode(ctx.w, ctx.r, ui.SignupForm("#modal-container"))
			return nil, err
		}

		// Redirect to the default recipes page for the user
		url := fmt.Sprintf("%s/g/%d/recipes", appconfig.Config.URL, groups[0].ID)
		ctx.w.Header().Set("HX-Redirect", url)
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
			return nil, err
		}

		// Get email user entered
		err = ctx.r.ParseForm()
		if err != nil {
			return nil, err
		}
		email := ctx.r.FormValue("email")

		desination := &types.Destination{
			ToAddresses: []string{email},
		}

		// Use users email to create a auth token
		// Auth token will allow use to verify their email
		token, err := h.CreateRegistrationToken(ctx.context(), email)
		if err != nil {
			return nil, err
		}

		// Construct the magic link to be used for verification
		magicLink := appconfig.Config.URL + "/auth/verify?token=" + token

		params := &sesv2.SendEmailInput{
			Content:              createLoginEmail(appconfig.AppName(), magicLink),
			ConfigurationSetName: new(string),
			Destination:          desination,
			FromEmailAddress:     &appconfig.Config.FromEmail,
		}

		client := sesv2.NewFromConfig(config)
		_, err = client.SendEmail(ctx.context(), params)
		if err != nil {
			slog.Error("Could not send email", "to", email, "err", err.Error())
			return nil, err
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
			return nil, nil
		}

		slog.Info("Verified email")

		user, err := h.GetUser(ctx.context(), email)
		if err != nil {
			// Assume account does not exist
			//return ui.CreateAccountPage(ui.PageProps{}), nil
			user, err = h.CreateAccount(ctx.context(), email)
			if err != nil {
				return nil, err
			}

		}

		session_token := generateSessionToken(32)
		loggedIn, err := h.Login(ctx.context(), user.ID, session_token)
		if err != nil {
			return nil, err
		}
		if !loggedIn {
			return nil, nil
		}

		var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		store.Options.HttpOnly = true

		session, _ := store.Get(ctx.r, "session")
		// Set some session values.
		session.Values["session_token"] = session_token
		err = session.Save(ctx.r, ctx.w)
		if err != nil {
			return nil, err
		}

		groups, err := h.GetUserGroups(ctx.context(), user.ID)
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("%s/g/%d/recipes", appconfig.Config.URL, groups[0].ID)
		http.Redirect(ctx.w, ctx.r, url, http.StatusSeeOther)
		slog.Info("Redirect verified user", "url", url)

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
	ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return node, nil
	})(w, r)
}
