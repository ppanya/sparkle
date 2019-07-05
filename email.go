package sparkle

type EmailSender interface {
	Send(to, from, content string) error
}
