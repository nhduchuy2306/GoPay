package notification

import "log"

type Service struct {
	mailService *EmailSender
}

func (s *Service) SendTransferSuccess(to string, content string) {
	if err := s.mailService.Send(to, content); err != nil {
		log.Fatal(err)
	}
}
