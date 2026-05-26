package notification

import (
	"log"
	"net/smtp"
)

type ConsoleService struct{}

func NewConsoleService() *ConsoleService {
	return &ConsoleService{}
}

func (s *ConsoleService) SendEmail(to, subject, body string) error {
	log.Printf("[EMAIL] To: %s | Subject: %s | Body: %s", to, subject, body)
	return nil
}

type SMTPService struct {
	cfg Config
}

func NewSMTPService(cfg Config) *SMTPService {
	return &SMTPService{cfg: cfg}
}

func (s *SMTPService) SendEmail(to, subject, body string) error {
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, msg)
}
