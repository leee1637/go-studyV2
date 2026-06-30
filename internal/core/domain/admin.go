package domain

type Admin struct {
	ID          int     `json"id"`
	FIO         string  `json"fio"`
	PhoneNumber *string `json"phone_number"`
}
