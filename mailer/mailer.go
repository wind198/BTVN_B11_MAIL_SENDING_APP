package mailer

import (
	"log"

	"example.com/rabbitmq/btvn_b11/scan_publish"
	u "example.com/rabbitmq/btvn_b11/utility"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	// "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	client *sendgrid.Client
}

func NewSendGridClient(api string) *SendGridClient {
	var client SendGridClient
	client.client = sendgrid.NewSendClient(api)
	return &client
}

func (s *SendGridClient) Send(email scan_publish.Email) error {
	from := mail.NewEmail("", email.From)
	to := mail.NewEmail("", email.To)
	subject := email.Subject
	plainContent := email.Content
	HTMLcontent := email.HTMLcontent
	message := mail.NewSingleEmail(from, subject, to, plainContent, HTMLcontent)
	response, err := s.client.Send(message)
	u.PrintError(err, "Error sending email")
	if err == nil {
		log.Println("Send message sucessfully", response, response.StatusCode)
		return nil
	} else {
		return err
	}

}
