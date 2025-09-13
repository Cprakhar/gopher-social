package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey string
	client *sendgrid.Client
}

func NewSendGrid(fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey: apiKey,
		client: client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)
	
	// template parsing and dynamic data handling would go here
	tmpl, err := template.ParseFS(FS, "templates/" + templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},

	})
	
	var retryErr error
	for i := range MaxRetries {
		response, retryErr := s.client.Send(message)
		if retryErr != nil {
			// exponential backoff could be implemented here
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts, error: %v", MaxRetries, retryErr)
}