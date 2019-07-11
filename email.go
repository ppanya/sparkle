package sparkle

import "fmt"

type EmailSender interface {
	Send(to, from, subject, content string, options ...EmailOption) error
}

type EmailOptions struct {
	CC  []string
	BCC []string
}

type EmailOption func(options *EmailOptions)

func WithCCOption(CC ...string) EmailOption {
	return func(options *EmailOptions) {
		options.CC = CC
	}
}

func WithBCCOption(BCC ...string) EmailOption {
	return func(options *EmailOptions) {
		options.BCC = BCC
	}
}

func ComposeEmailOptions(options ...EmailOption) *EmailOptions {
	var opts = EmailOptions{
		BCC: []string{},
		CC:  []string{},
	}

	for _, option := range options {
		option(&opts)
	}
	return &opts
}

type ConsoleEmailSender struct {
	Inbox [][]string
}

func (c *ConsoleEmailSender) Send(to, from, subject, content string, options ...EmailOption) error {
	option := ComposeEmailOptions(options...)
	message := []string{to, from, subject, content}

	for _, bcc := range option.BCC {
		message = append(message, bcc)
	}

	for _, cc := range option.CC {
		message = append(message, cc)
	}

	c.Inbox = append(c.Inbox, message)
	fmt.Printf("email sent!!\nfrom:%s -> to:%s", from, to)
	fmt.Printf("subject:%s\n", subject)
	fmt.Printf("content:%s\n", content)
	fmt.Printf("BCC:%v\n", option.BCC)
	fmt.Printf("CC:%v\n", option.CC)
	return nil

}
