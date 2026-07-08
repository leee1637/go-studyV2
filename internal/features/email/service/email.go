package service_email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	EmailConfig EmailConfig
}

type EmailConfig struct {
	Host     string
	Port     int16
	Username string
	Password string
	From     string
	BaseURL  string
}

func NewEmailConfig(em EmailConfig) *EmailService {
	return &EmailService{
		EmailConfig: em,
	}
}

func (e *EmailService) SendVerificationLink(to, fio, token string) error {
	verifLink := fmt.Sprintf("%s/api/auth/verify/%s", e.EmailConfig.BaseURL, token)

	htmlBody := fmt.Sprintf(`
        <h2>Подтверждение регистрации</h2>
        <p>Здравствуйте, %s!</p>
        <p>Для завершения регистрации перейдите по ссылке:</p>
        <p><a href="%s">%s</a></p>
        <p>Ссылка действительна в течение 24 часов.</p>
        <p>Если вы не регистрировались на нашем сайте, просто игнорируйте это письмо.</p>
    `, fio, verifLink, verifLink)

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.EmailConfig.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Подтверждение регистрации")
	msg.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(e.EmailConfig.Host, int(e.EmailConfig.Port), e.EmailConfig.Username, e.EmailConfig.Password)

	err := dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("Ошибка отправки сообщения: %w", err)
	}

	return nil
}
