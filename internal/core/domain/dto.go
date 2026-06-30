package domain

type SignUpDTO struct {
	ID          int
	Login       string
	Password    string
	Role        Role
	FIO         string
	GroupName   []string
	PhoneNumber *string
}
