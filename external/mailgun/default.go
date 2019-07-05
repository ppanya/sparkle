package mailgun

import "fmt"

type MailGunEmailSender struct{}

func (m *MailGunEmailSender) Send(to, from, content string) error {

	fmt.Println(to)
	fmt.Println(content)

	return nil
}
