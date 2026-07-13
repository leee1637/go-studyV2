package domain

type SignUpDTO struct {
	ID          int
	Email       string
	Password    string
	Role        Role
	FIO         string
	GroupName   []string
	PhoneNumber *string
}

type UpdateStudentDTO struct {
	FIO         string  `json:"fio"`
	GroupName   string  `json:"group"`
	PhoneNumber *string `json:"phone_number"`
}

type UpdateTeacherDTO struct {
	FIO         string  `json:"fio"`
	PhoneNumber *string `json:"phone_number"`
}
