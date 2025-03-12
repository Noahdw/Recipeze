package http

import (
	//"context"

	"fmt"
	"log/slog"
	"net/http"
	"recipeze/appconfig"
	"recipeze/ui"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsc "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
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

	r.Post("/auth/magic-link", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		email := r.FormValue("email")
		config, err := awsc.LoadDefaultConfig(r.Context(),
			awsc.WithRegion("us-east-2"),
		)
		if err != nil {
			// TODO
			slog.Error("Could not load email config")
			return nil, err
		}

		desination := &types.Destination{
			ToAddresses: []string{email},
		}

		magicLink := appconfig.Config.URL + "/recipes"

		params := &sesv2.SendEmailInput{
			Content:              createLoginEmail("App", magicLink),
			ConfigurationSetName: new(string),
			Destination:          desination,
			FromEmailAddress:     &appconfig.Config.FromEmail,
		}

		client := sesv2.NewFromConfig(config)
		_, err = client.SendEmail(r.Context(), params)
		if err != nil {
			slog.Error("Could not send email", "to", email, "err", err.Error())
			return nil, err
		}

		slog.Info("sending magic link", "email", email)
		return ui.SignupForm("#modal-container"), nil
	}))
}

func createLoginEmail(appName, magicLink string) *types.EmailContent {
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login to %s</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f9f9f9;
            border-radius: 5px;
            padding: 20px;
        }
        .button {
            display: inline-block;
            background-color: #007bff;
            color: white;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            margin-top: 30px;
            font-size: 12px;
            color: #666666;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Login to %s</h2>
        <p>Hello,</p>
        <p>Click the button below to securely log in to your account. This link will expire in 15 minutes.</p>
        <!-- Email clients often have better support for table-based buttons -->
        <table cellpadding="0" cellspacing="0" border="0" style="margin: 20px 0;">
            <tr>
                <td align="center" bgcolor="#007bff" style="border-radius: 5px;">
                    <a href="%s" target="_blank" style="display: inline-block; padding: 10px 20px; font-size: 16px; color: white; text-decoration: none; border-radius: 5px; font-family: Arial, sans-serif;">Login Now</a>
                </td>
            </tr>
        </table>
        <p>If you didn't request this login link, you can safely ignore this email.</p>
        <p>If the button above doesn't work, copy and paste this URL into your browser:</p>
        <p>%s</p>
    </div>
    <div class="footer">
        <p>This is an automated message from %s. Please do not reply to this email.</p>
    </div>
</body>
</html>
`, appName, appName, magicLink, magicLink, appName)

	textBody := fmt.Sprintf(`
Login to %s

Hello,

Click the link below to securely log in to your account. This link will expire in 15 minutes.

%s

If you didn't request this login link, you can safely ignore this email.

This is an automated message from %s. Please do not reply to this email.
`, appName, magicLink, appName)

	emailContent := &types.EmailContent{
		Simple: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(textBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(fmt.Sprintf("Login to %s", appName)),
			},
		},
	}
	return emailContent
}
