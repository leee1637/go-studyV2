package transport_teacher

import service_teacher "study/internal/features/teachers/service"

type TeacherHandler struct {
	TeacherService *service_teacher.TeacherService
}

func NewTeacherHandler(TeacherService *service_teacher.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		TeacherService: TeacherService,
	}
}
