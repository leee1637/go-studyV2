package http_transport

import (
	service_registration "study/internal/features/registration/service"
	"study/internal/features/users/service"
)

type AuthHandler struct {
	authService *service.AuthService
	regService  *service_registration.RegistrationService
}

func NewAuthHandler(authService *service.AuthService, regService *service_registration.RegistrationService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		regService:  regService,
	}
}
