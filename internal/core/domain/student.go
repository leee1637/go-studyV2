package domain

type Student struct {
	ID          int     `json:"id"`
	FIO         string  `json:"fio"`
	GroupName   string  `json:"group"`
	PhoneNumber *string `json:"phone_number"`
}
