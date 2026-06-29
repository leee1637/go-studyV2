package domain

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleTeacher Role = "TEACHER"
	RoleStudent Role = "STUDENT"
)
