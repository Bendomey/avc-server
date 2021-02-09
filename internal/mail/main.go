package mail

import (
	"context"
	"time"

	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/mailgun/mailgun-go/v4"
)

var domain, apiKey, sender string

func init() {
	domain = utils.MustGet("MAILGUN_DOMAIN")
	apiKey = utils.MustGet("MAILGUN_API_KEY")
	sender = utils.MustGet("EMAIL_SENDER")
}

// MailingService interface
type MailingService interface {
	SendTransactionalMail(ctx context.Context, subject string, body string, recipient string) error
}

// Mail take obj of Mail
type Mail struct {
	mg *mailgun.MailgunImpl
}

// NewMailingSvc exposed the ORM to the user functions in the module
func NewMailingSvc() MailingService {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, apiKey)
	return &Mail{mg}
}

//SendTransactionalMail sends a transaction mail
func (mg *Mail) SendTransactionalMail(ctx context.Context, subject string, body string, recipient string) error {

	// The message object allows you to add attachments and Bcc recipients
	message := mg.mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err := mg.mg.Send(ctx, message)

	if err != nil {
		return err
	}
	return nil
}
