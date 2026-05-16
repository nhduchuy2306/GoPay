package notification

import "log"

type Service struct {
	emailSender *EmailSender
}

func NewService(emailSender *EmailSender) *Service {
	return &Service{emailSender: emailSender}
}

func (s *Service) SendWelcomeEmail(to string) error {
	if err := s.emailSender.SendWelcome(to); err != nil {
		log.Println("welcome email error:", err)
		return err
	}
	return nil
}

func (s *Service) SendTransferSuccessEmail(to string, data TransferSuccessEmailData) error {
	if err := s.emailSender.SendTransferSuccess(to, data); err != nil {
		log.Println("transfer success email error:", err)
		return err
	}
	return nil
}

var _ = NewService
