package http

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

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
