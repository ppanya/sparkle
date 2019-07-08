package sparkle

import "fmt"

type EmailSender interface {
	Send(to, from, content string) error
}

type ConsoleEmailSender struct {
	Inbox [][]string
}

func (c *ConsoleEmailSender) Send(to, from, content string) error {
	c.Inbox = append(c.Inbox, []string{to, from, content})
	fmt.Printf("email sent!\n%s\n%s\n%s\n", to, from, content)
	return nil

}
