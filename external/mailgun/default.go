package mailgun

import (
	"context"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/octofoxio/sparkle"
	"time"
)

type mailGunEmailSender struct {
	client mailgun.Mailgun
}

func (m *mailGunEmailSender) Send(to, from, subject, content string, options ...sparkle.EmailOption) error {
	option := sparkle.ComposeEmailOptions(options...)
	message := m.client.NewMessage(from, subject, "", to)
	message.SetHtml(content)
	for _, bcc := range option.BCC {
		message.AddBCC(bcc)
	}

	for _, cc := range option.CC {
		message.AddCC(cc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := m.client.Send(ctx, message)
	return err
}

func NewMailGunEmailSender(domain, apiKey string) *mailGunEmailSender {
	return &mailGunEmailSender{
		client: mailgun.NewMailgun(domain, apiKey),
	}
}
