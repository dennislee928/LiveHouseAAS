package notification

type Service interface {
	SendEmail(to, subject, body string) error
}

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}
