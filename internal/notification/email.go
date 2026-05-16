package notification

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopay/internal/mail"
)

type EmailSender struct {
	mailer      *mail.Service
	templateDir string
}

type WelcomeEmailData struct {
	Name  string
	Email string
}

type TransferSuccessEmailData struct {
	Recipient     string
	Amount        int64
	Currency      string
	TransactionID string
	FromWalletID  string
	ToWalletID    string
	Description   string
}

func NewEmailSender(mailer *mail.Service) *EmailSender {
	return &EmailSender{
		mailer:      mailer,
		templateDir: filepath.Join("internal", "templates", "email"),
	}
}

func (e *EmailSender) SendWelcome(to string) error {
	data := WelcomeEmailData{
		Name:  friendlyNameFromEmail(to),
		Email: to,
	}
	body, err := mail.ParseTemplate(filepath.Join(e.templateDir, "welcome.html"), data)
	if err != nil {
		return err
	}
	return e.mailer.SendHTML(to, "Welcome to Go pay", body)
}

func (e *EmailSender) SendTransferSuccess(to string, data TransferSuccessEmailData) error {
	data.Recipient = friendlyNameFromEmail(to)
	body, err := mail.ParseTemplate(filepath.Join(e.templateDir, "transfer_success.html"), data)
	if err != nil {
		return err
	}
	return e.mailer.SendHTML(to, "Transfer received successfully", body)
}

func friendlyNameFromEmail(email string) string {
	local := strings.Split(email, "@")[0]
	if local == "" {
		return "there"
	}
	parts := strings.FieldsFunc(local, func(r rune) bool { return r == '.' || r == '_' || r == '-' })
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s", local)
	}
	return strings.Join(parts, " ")
}
