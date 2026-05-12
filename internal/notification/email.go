package notification

import "log"

type EmailSender struct {
}

func (e *EmailSender) Send(to string, message string) error {
	log.Println("Send email to", to, ":", message)
	return nil
}
