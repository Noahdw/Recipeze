package handler

import (
	"log/slog"
	"net/http"
	"recipeze/appconfig"

	"recipeze/ui"

	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
)

func (h *handler) RouteHome(r chi.Router) {
	r.Get("/", h.adapt(func(ctx requestContext) (Node, error) {
		return ui.HomePage(ui.PageProps{}), nil
	}))

	r.Get("/login", h.adapt(func(ctx requestContext) (Node, error) {
		return ui.SignupForm("#modal-container"), nil
	}))

	r.Post("/auth/magic-link", h.sendMagicLinkToEmail())

	r.Get("/auth/verify", h.authenticateByMagicLink())
}

func (h *handler) authenticateByMagicLink() http.HandlerFunc {
	return h.adapt(func(ctx requestContext) (Node, error) {
		token := ctx.queryParam("token")
		email, err := h.VerifyRegistrationToken(ctx.context(), token, ctx.r)
		if err != nil {
			slog.Error("Could not verify email", "err", err.Error())
			url := appconfig.Config.URL
			http.Redirect(ctx.w, ctx.r, url, http.StatusTemporaryRedirect)
			return nil, nil
		}

		slog.Info("Verified email")

		_, err = h.GetUser(ctx.context(), email)
		if err != nil {
			// Assume account does not exist
			return ui.CreateAccountPage(ui.PageProps{}), nil
		}

		url := appconfig.Config.URL + "/recipes"
		http.Redirect(ctx.w, ctx.r, url, http.StatusTemporaryRedirect)

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
