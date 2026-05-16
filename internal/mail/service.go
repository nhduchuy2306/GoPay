package mail

import "gopkg.in/gomail.v2"

type Service struct {
	Host string
	Port int
	From string
}

func NewService(host string, port int) *Service {
	return &Service{Host: host, Port: port, From: "no-reply@gopay.local"}
}

func (s *Service) SendHTML(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.Host, s.Port, "", "")

	return d.DialAndSend(m)
}
