package dispatcher

import (
	service_email "study/internal/features/email/service"
)

type DirectDispatcher struct {
	emailService *service_email.EmailService
}

func NewDirectDispatcher(emailService *service_email.EmailService) *DirectDispatcher {
	return &DirectDispatcher{emailService: emailService}
}

func (d *DirectDispatcher) Send(to, fio, token string) error {
	return d.emailService.SendVerificationLink(to, fio, token)
}
