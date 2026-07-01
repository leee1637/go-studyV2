package http_student

import student_service "study/internal/features/students/service"

type StudentHandler struct {
	StudentService *student_service.StudentService
}

func NewAuthService(s student_service.StudentService) *StudentHandler {
	return &StudentHandler{
		StudentService: &s,
	}
}
